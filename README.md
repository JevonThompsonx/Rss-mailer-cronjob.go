# cronJobs
cron jobs for one of my servers that would be useful for me. Right now it'll be just rss emails but who knows 

## prereqs 
- go is setup and ready to use: [link to setup](https://docs.google.com/document/d/1QkiZEAUWcW6f5Ep73_DzvOcPKxTvfwbZX0M-cO0Id90/edit?usp=sharing)

## Use
1. git clone repo
2. cd into repo `cd cronJobs`
3. `go mod init go-rss-mailer`
4. Add a email and password login to a .env file, use app passwords if possible. Google provides app passwords. It might be a good idea to create a totally new email just for this sort of automated messaging
  - `nano .env`
  - ```
    # .env file for SMTP Credentials
    SMTP_USERNAME="your-email@gmail.com"
    SMTP_PASSWORD="your-gmail-app-password"
    ```
5. Go get needed dotenv repo: `# In ~/go-rss-mailer
go get github.com/joho/godotenv`
6. tidy: `go mod tidy`
5. Build `go build .`
6. Test `./go-rss-mailer`
7. Add combiled task to cron job `crontab -e`
  - Fetch RSS feeds every hour at 15 past the hour: make sure path is correct
`15 * * * * cd /home/jevonx/go-rss-mailer && ./go-rss-mailer >> /home/jevonx/go-rss-mailer/mailer.log 2>&1`
