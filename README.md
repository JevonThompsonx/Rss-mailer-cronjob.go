# RSS Mailer Cron Job

This repository contains a Go application designed to be run as a cron job. It fetches articles from a list of RSS/Atom feeds, checks for new entries since the last run, and emails a digest of the new articles to a specified address.

## Prerequisites

* A Linux server (tested on Debian 12).
* Go (version 1.20+) is installed and configured in your `$PATH`. You can find installation instructions on the [official Go website](https://go.dev/doc/install) or [here]([link to setup](https://docs.google.com/document/d/1QkiZEAUWcW6f5Ep73_DzvOcPKxTvfwbZX0M-cO0Id90/edit?usp=sharing).
* An email account that can be used for sending mail via SMTP. If using Gmail, you must set up a 2-Factor Authentication and create an [App Password](https://support.google.com/accounts/answer/185833).

## Installation and Setup

Follow these steps to get the service running.

## Use
1. git clone repo
2. cd into repo `cd cronJobs`
3. Add a email and password login to a .env file, use app passwords if possible. Google provides app passwords. It might be a good idea to create a totally new email just for this sort of automated messaging
   - `cp .env.example .env && nano .env`
4. tidy/get required modules: `go mod tidy`
5. Build `go build .`
6. Test `./go-rss-mailer`
7. Add combiled task to cron job `crontab -e`
  - Fetch RSS feeds every hour at 15 past the hour: make sure path is correct
`15 * * * * cd /home/jevonx/go-rss-mailer && ./go-rss-mailer >> /home/jevonx/go-rss-mailer/mailer.log 2>&1`

### Talking to myself...

This project can be used on a server or on my most used machine. I can clone it, install using `go install .` then add it as a startup command to my shell. 
For now, it might be most useful on my most used machine since I don't yet have a extremely low powered server I'd want to add this to. I could just get a pi tho
