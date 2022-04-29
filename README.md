# markdown-tweet-scheduler
Schedule tweets from markdown files in your repo, posted to twitter via github actions. 

Now supports .txt as well as .md files!


## Why
- View, edit, and post your tweets without logging into twitter 
- Keep source controlled backups of your tweets
- Free and open source


## Setup 
1. Clone this repo and push it to your own private repo. 
2. Get your credentials by creating a twitter app (https://developer.twitter.com/apps)
3. Add your twitter credentials to the repository's secrets (https://docs.github.com/en/actions/reference/encrypted-secrets)
   - API_KEY (known as consumer_key in twitter API)
   - API_SECRET_KEY (known as consumer_secret in twitter API)
   - ACCESS_TOKEN
   - ACCESS_TOKEN_SECRET


## Scheduling Tweets By Date
1. Create a markdown file in the `./tweets/` folder with a future date in either the `YYYY-Mon-DD` or `Month dd, YYYY` format, for example, `2021-Sep-05.md` or `September 5, 2021.md`, and write the content of your tweet in it.
2. Commit the file(s) and push to the remote repo. When the daily action runs on the specified date, the tweet should be posted.


## Configuration
1. By default, posts tweets around 7:02AM PT (2:02 UTC). To change the time of day tweets are posted:
   - Set the `cron` section of `.github/workflows/go.yml` to the time you want the tweet to post  https://cron.help/

2. Changing the tweet directory:
   - Change the `FILE_PATH` environment variable in `.github/workflows/go.yml`


## Running locally
1. rename `.env-SAMPLE` to `.env` and fill in your twitter credentials
   - consumer_key == API_KEY
   - consumer_secret == API_SECRET_KEY
2. In your terminal, from the root directory: `go run .`


## Notes
1. Fails silently on bad credentials, make sure you set those correctly.
2. Only allows one tweet per day by design. If requested, this can be modified to allow tweets by the minute or hour. 
3. Tweets will not be posted exactly at the cron time set in `go.yml`. In my experience in can be 5-10 minutes late. If you need minute precision, run this script locally on a cronjob.

