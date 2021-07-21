package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func Tweet(content string) {
	LoadDotEnv()

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", os.Getenv("API_KEY"), "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", os.Getenv("API_SECRET_KEY"), "Twitter Consumer Secret")
	accessToken := flags.String("access-token", os.Getenv("ACCESS_TOKEN"), "Twitter Access Token")
	accessSecret := flags.String("access-secret", os.Getenv("ACCESS_TOKEN_SECRET"), "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
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
	user, _, _ := client.Accounts.VerifyCredentials(verifyParams)
	fmt.Printf("User's ACCOUNT:\n%+v\n", user)

	tweet, _, _ := client.Statuses.Update(content, nil)
	fmt.Printf("Posted Tweet\n%v\n", tweet)
}

func ReadFile(file_name string) string {
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, _ := ioutil.ReadAll(file)
	return string(b)
}

func main() {
	LoadDotEnv()
	current_time := time.Now()
	formatted_time := current_time.Format("2006-Jan-02")
	path := os.Getenv("FILE_PATH")
	fmt.Println(path)
	todays_file_name := path + formatted_time + ".md"
	fmt.Println(todays_file_name)
	content := ReadFile(todays_file_name)

	Tweet(content)
}
