package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
)

type Config struct {
	SMTPServer string   `json:"smtp_server"`
	SMTPPort   int      `json:"smtp_port"`
	ToEmail    string   `json:"to_email"`
	FromEmail  string   `json:"from_email"`
	FeedURLs   []string `json:"feed_urls"`
}

type State struct {
	SentGUIDs map[string]bool `json:"sent_guids"`
}

const (
	configFile = "config.json"
	stateFile  = "sent_items.json"
)

func main() {
	// 1. Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables.")
	}

	// 2. Load Configuration from JSON
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 3. Load State (history of sent items)
	state, err := loadState(stateFile)
	if err != nil {
		log.Printf("Could not load state file, assuming new start: %v", err)
	}

	fp := gofeed.NewParser()
	var newItems []*gofeed.Item

	// 4. Fetch all feeds and collect new items
	fmt.Println("Checking for new articles...")
	for _, url := range config.FeedURLs {
		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Printf("Error fetching feed %s: %v", url, err)
			continue
		}

		for _, item := range feed.Items {
			if _, found := state.SentGUIDs[item.GUID]; !found {
				newItems = append(newItems, item)
			}
		}
	}

	if len(newItems) == 0 {
		fmt.Println("No new articles found.")
		return
	}

	fmt.Printf("Found %d new articles. Preparing email...\n", len(newItems))

	// --- CHANGE: Randomize the order of new articles ---
	// Seed the random number generator to ensure different order each time.
	rand.Seed(time.Now().UnixNano())
	// Shuffle the slice of new items in place.
	rand.Shuffle(len(newItems), func(i, j int) {
		newItems[i], newItems[j] = newItems[j], newItems[i]
	})
	// --- END CHANGE ---

	// 5. Get SMTP credentials from environment
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	if smtpUser == "" || smtpPass == "" {
		log.Fatal("Error: SMTP_USERNAME or SMTP_PASSWORD environment variables not set.")
	}

	// 6. Build and Send the Email
	err = sendEmail(config, smtpUser, smtpPass, newItems)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	// 7. Update State with newly sent items
	for _, item := range newItems {
		state.SentGUIDs[item.GUID] = true
	}
	err = saveState(stateFile, state)
	if err != nil {
		log.Fatalf("Failed to save state: %v", err)
	}

	fmt.Printf("Successfully sent %d new articles and updated state.\n", len(newItems))
}

func loadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	return &config, err
}

func loadState(path string) (*State, error) {
	file, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{SentGUIDs: make(map[string]bool)}, nil
	}
	if err != nil {
		return nil, err
	}
	var state State
	err = json.Unmarshal(file, &state)
	return &state, err
}

func saveState(path string, state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func sendEmail(config *Config, user, pass string, items []*gofeed.Item) error {
	var body strings.Builder
	body.WriteString("<h1>Your News Digest</h1>")
	for _, item := range items {
		// This part remains unchanged, but now processes items in random order.
		var publishedDate string
		if item.PublishedParsed != nil {
			publishedDate = item.PublishedParsed.Format(time.RFC822)
		} else {
			publishedDate = "N/A"
		}
		body.WriteString(fmt.Sprintf(
			`
			<hr>
			<h3><a href="%s">%s</a></h3>
			<p><i>Published: %s</i></p>
			<p>%s</p>
			`,
			template.HTMLEscapeString(item.Link),
			template.HTMLEscapeString(item.Title),
			publishedDate,
			item.Description,
		))
	}

	e := email.NewEmail()
	e.From = config.FromEmail
	e.To = []string{config.ToEmail}
	e.Subject = fmt.Sprintf("News Digest: %d New Articles", len(items))
	e.HTML = []byte(body.String())

	addr := fmt.Sprintf("%s:%d", config.SMTPServer, config.SMTPPort)
	auth := smtp.PlainAuth("", user, pass, config.SMTPServer)
	return e.Send(addr, auth)
}
