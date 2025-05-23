# cronJobs
cron jobs for one of my servers that would be useful for me. Right now it'll be just rss emails but who knows 

## Use
1. git clone repo
2. cd into repo `cd cronJobs`
3. Add a email and password login to a .env file, use app passwords if possible. Google provides app passwords. It might be a good idea to create a totally new email just for this sort of automated messaging
  - `nano .env`
  - ```
    # .env file for SMTP Credentials
    SMTP_USERNAME="your-email@gmail.com"
    SMTP_PASSWORD="your-gmail-app-password"
    ```
5. Build `go build .`
6. Add combiled task to cron job `/etc/cron`
  - 
