# markdown-tweet-scheduler
Schedule daily tweets from markdown files in your source controlled repository. 

## Setup 
1. Clone this repo
2. Get your twitter credentials (https://developer.twitter.com/en)
3. Add your twitter credentials to the repository's secrets (https://docs.github.com/en/actions/reference/encrypted-secrets)
   API_KEY
   API_SECRET_KEY
   ACCESS_TOKEN
   ACCESS_TOKEN_SECRET
4. Create a markdown file in the `./tweets/` folder with tomorrow's date in `YYYY-Mon-DD` format, for example, `2021-Aug-05.md`, and write the content of your tweet in it.

## Why
- View, edit, and post your tweets without logging into twitter manually
- Keep source controlled backups of your tweets
- Free and open source

## Configuration
1. By default, posts tweets around 7:02AM PT (2:02 UTC). To change the time of day tweets are posted:
  Set the `cron` section of `.github/workflows/go.yml` to the time you want the tweet to post (note: it will not run exactly at this time) https://cron.help/

2. Changing the tweet directory:
   Change the `FILE_PATH` environment variable in `.github/workflows/go.yml` 

## Running locally
1. rename `.env-SAMPLE` to `.env` and fill in your twitter credentials
   - consumer_key == API_KEY
   - consumer_secret == API_SECRET_KEY
2. In your terminal (if it's the first time running): `go mod init main.go`
3. In your terminal: `go run main.go`

## Notes
1. Fails silently on bad credentials, make sure you set those correctly.
2. Only allows one tweet per day by design. If requested, this can be modified to allow tweets by the minute or hour. 
3. I suggest moving posted tweets to a `posted/` subdirectory under `tweets/`, but it's completely optional.
