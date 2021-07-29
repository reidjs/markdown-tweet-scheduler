# markdown-tweet-scheduler
Schedule daily tweets from markdown files in your repo, posted to twitter via github actions. 

## Setup 
1. Fork this repo
2. Get your twitter credentials (https://developer.twitter.com/en)
3. Add your twitter credentials to the repository's secrets (https://docs.github.com/en/actions/reference/encrypted-secrets)
   - API_KEY (known as consumer_key in twitter API)
   - API_SECRET_KEY (known as consumer_secret in twitter API)
   - ACCESS_TOKEN
   - ACCESS_TOKEN_SECRET

## Scheduling Tweets
**By Queue**
1. Create a markdown file in the `./tweets/` folder with a name starting with `q-` for example, `q-1.md` or `q-84.md`
2. Commit the file(s) and push to the remote repo. When the daily action runs, the first queued tweet in alphanumeric order should be posted. 

**By Date**
1. Create a markdown file in the `./tweets/` folder with tomorrow's date in `YYYY-Mon-DD` format, for example, `2021-Sep-05.md`, and write the content of your tweet in it.
2. Commit the file(s) and push to the remote repo. When the daily action runs on the specified date, the tweet should be posted.

**Both**
If you have both date specified and queue specified posts in the `tweets` folder, date specified tweets will take precedence. I.e., the queued tweet(s) will wait until the next day in which there isn't a date specified tweet to post. Only one tweet may be posted by day, by design. 

## Why
- View, edit, and post your tweets without logging into twitter manually
- Keep source controlled backups of your tweets
- Free and open source

## Configuration
1. By default, posts tweets around 7:02AM PT (2:02 UTC). To change the time of day tweets are posted:
   - Set the `cron` section of `.github/workflows/go.yml` to the time you want the tweet to post  https://cron.help/

2. Changing the tweet directory:
   - Change the `FILE_PATH` environment variable in `.github/workflows/go.yml`

## Running locally
1. rename `.env-SAMPLE` to `.env` and fill in your twitter credentials
   - consumer_key == API_KEY
   - consumer_secret == API_SECRET_KEY
2. In your terminal: `go run main.go`

## Notes
1. Fails silently on bad credentials, make sure you set those correctly.
2. Only allows one tweet per day by design. If requested, this can be modified to allow tweets by the minute or hour. 
3. I suggest moving posted tweets that have already been posted to a `posted/` subdirectory under `tweets/`.
4. Tweets will not be posted exactly at the cron time set in `go.yml` because of how github actions work. If you need minute precision, run this script locally on a cronjob.
