package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file")
	}
}

func Tweet(content string) error {
	LoadDotEnv()

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", os.Getenv("API_KEY"), "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", os.Getenv("API_SECRET_KEY"), "Twitter Consumer Secret")
	accessToken := flags.String("access-token", os.Getenv("ACCESS_TOKEN"), "Twitter Access Token")
	accessSecret := flags.String("access-secret", os.Getenv("ACCESS_TOKEN_SECRET"), "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		// log.Fatal("Consumer key/secret and Access token/secret required")
		return errors.New("Bad credentials")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}
	user, _, err1 := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's ACCOUNT:\n%+v\n", user)
	if err1 != nil {
		return err1
	}

	tweet, _, err2 := client.Statuses.Update(content, nil)
	fmt.Printf("Posted Tweet\n%v\n", tweet)
	return err2
}

func ReadFile(file_name string) (string, error) {
	file, err := os.Open(file_name)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, _ := ioutil.ReadAll(file)
	return string(b), nil
}

func ScheduledTweet() (string, string, error) {
	LoadDotEnv()
	current_time := time.Now()
	formatted_time := current_time.Format("2006-Jan-02")
	path := os.Getenv("FILE_PATH")
	fmt.Println(path)
	todays_file_name := path + formatted_time + ".md"
	fmt.Println(todays_file_name)
	content, err := ReadFile(todays_file_name)
	if err != nil {
		return "", "", err
	}

	return content, todays_file_name, nil
}

func IsQueueNameFormat(filename string) bool {
	if string(filename[0]) == "q" && string(filename[1]) == "-" {
		return true
	}
	return false
}

func QueuedTweet() (string, string, error) {
	LoadDotEnv()
	path := os.Getenv("FILE_PATH")
	files, dir_read_error := ioutil.ReadDir(path)
	if dir_read_error != nil {
		return "", "", dir_read_error
	}
	filenames := []string{}

	for _, f := range files {
		// fmt.Println(f.Name())
		// if name fits the format q-#, append it
		if IsQueueNameFormat(f.Name()) {
			filenames = append(filenames, f.Name())
		}
	}

	sort.Strings(filenames)
	if len(filenames) == 0 {
		return "", "", errors.New("No queued files found")
	}
	queued_file_name := filenames[0]

	queued_file_path := path + queued_file_name

	content, read_error := ReadFile(queued_file_path)
	if read_error != nil {
		return "", "", read_error
	}

	current_time := time.Now()
	formatted_time := current_time.Format("2006-Jan-02")
	new_filename_path := path + "attempted_" + formatted_time + "_" + queued_file_name
	rename_error := os.Rename(queued_file_path, new_filename_path)
	if rename_error != nil {
		return "", "", rename_error
	}

	// return the content of the file
	return content, new_filename_path, nil
}

func main() {
	scheduled_content, scheduled_tweet_filename, scheduled_tweet_error := ScheduledTweet()
	path := os.Getenv("FILE_PATH")

	if scheduled_tweet_error != nil {
		fmt.Println(scheduled_tweet_error)
		// If there's not a scheduled tweet today, then try a queued tweet
		queued_content, queued_tweet_filename, queued_tweet_error := QueuedTweet()
		if queued_tweet_error != nil {
			fmt.Println(queued_tweet_error)
		} else {
			fmt.Println("Posting", queued_content)
			post_failure := Tweet(queued_content)
			if post_failure != nil {
				fmt.Println("todo: change filename to failed_", queued_tweet_filename)
				// TODO: check if this is correct naming
				os.Rename(path+queued_tweet_filename, path+"failed_"+queued_tweet_filename[9:])
			} else {
				// TODO
				fmt.Println("todo: change filename to posted_", queued_tweet_filename)
			}
		}
	} else {
		post_failure := Tweet(scheduled_content)
		if post_failure != nil {
			// TODO
			fmt.Println("todo: change filename to failed_", scheduled_tweet_filename)
		} else {
			// TODO
			fmt.Println("todo: change filename to posted_", scheduled_tweet_filename)
		}
	}
}
