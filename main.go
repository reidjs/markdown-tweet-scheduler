package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

var (
	ErrNoFiles        = errors.New("no tweet files found")
	ErrBadCredentials = errors.New("bad credentials")
)

type Configuration struct {
	APIKey            string
	APISecret         string
	AccessToken       string
	AccessTokenSecret string
	FilePath          string
}

func main() {
	config, err := LoadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	content, err := config.ReadScheduledTweet()
	if err != nil {
		if errors.Is(err, ErrNoFiles) {
			log.Println("No files to tweet. Sleep.")
			os.Exit(0)
		}

		log.Fatal(err)
	}

	if err := config.PostTweet(content); err != nil {
		log.Fatalf("Error posting to Twitter: %v", err)
	}
}

func LoadConfiguration() (*Configuration, error) {
	if err := godotenv.Load(); err != nil || !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file %w", err)
	}

	config := Configuration{
		APIKey:            os.Getenv("API_KEY"),
		APISecret:         os.Getenv("API_SECRET_KEY"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		FilePath:          os.Getenv("FILE_PATH"),
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("%w: twitter consumer key is empty", ErrBadCredentials)
	}

	if config.APISecret == "" {
		return nil, fmt.Errorf("%w: twitter consumer secret is empty", ErrBadCredentials)
	}

	if config.AccessToken == "" {
		return nil, fmt.Errorf("%w: twitter access token is empty", ErrBadCredentials)
	}

	if config.AccessTokenSecret == "" {
		return nil, fmt.Errorf("%w: twitter access token secret is empty", ErrBadCredentials)
	}

	if config.FilePath == "" {
		return nil, fmt.Errorf("%w: file path is empty", ErrNoFiles)
	}

	return &config, nil
}

func (c *Configuration) ReadScheduledTweet() (string, error) {
	var (
		currentTime   = time.Now()
		path          = c.FilePath
		isoDate       = currentTime.Format("2006-Jan-02")
		fullDate      = currentTime.Format("January 2, 2006")
		possibleFiles = []string{
			fmt.Sprintf("%s%s.md", path, isoDate),
			fmt.Sprintf("%s%s.txt", path, isoDate),
			fmt.Sprintf("%s%s.md", path, fullDate),
			fmt.Sprintf("%s%s.txt", path, fullDate),
		}
		existingFiles = make([]string, 0, 1)
	)

	log.Println("File path: ", path)

	for _, file := range possibleFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		existingFiles = append(existingFiles, file)
	}

	if len(existingFiles) == 0 {
		return "", ErrNoFiles
	}

	firstFilepath := existingFiles[0]

	content, err := ReadFile(firstFilepath)
	if err != nil {
		return "", fmt.Errorf("reading from %s: %w", firstFilepath, err)
	}

	log.Println("Attempting to post content from: ", firstFilepath)

	return content, nil
}

func ReadFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", filename, err)
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)

	return string(b), nil
}

func (c *Configuration) PostTweet(content string) error {
	var (
		config = oauth1.NewConfig(c.APIKey, c.APISecret)
		token  = oauth1.NewToken(c.AccessToken, c.AccessTokenSecret)

		// OAuth1 http.Client will automatically authorize Requests
		httpClient = config.Client(oauth1.NoContext, token)

		// Twitter client
		client = twitter.NewClient(httpClient)

		// Verify Credentials
		verifyParams = &twitter.AccountVerifyParams{
			SkipStatus:   twitter.Bool(true),
			IncludeEmail: twitter.Bool(true),
		}
	)

	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return fmt.Errorf("verify credentials: %w", err)
	}

	if _, _, err := client.Statuses.Update(content, nil); err != nil {
		return fmt.Errorf("post tweet: %w", err)
	}

	return nil
}

// TODO: either fix or remove queue system.
func IsQueueNameFormat(filename string) bool {
	// if name fits the format q-#, return true
	return strings.HasPrefix(filename, "q-")
}

// TODO: either fix or remove queue system.
func (c *Configuration) QueuedTweet() (string, error) {
	files, err := ioutil.ReadDir(c.FilePath)
	if err != nil {
		return "", fmt.Errorf("read dir %s: %w", c.FilePath, err)
	}

	var filenames []string

	for _, f := range files {
		if IsQueueNameFormat(f.Name()) {
			filenames = append(filenames, f.Name())
		}
	}

	sort.Strings(filenames)

	if len(filenames) == 0 {
		return "", ErrNoFiles
	}

	var (
		queuedFilename = filenames[0]
		queuedFilepath = fmt.Sprintf("%s%s", c.FilePath, queuedFilename)
	)

	content, err := ReadFile(queuedFilepath)
	if err != nil {
		return "", fmt.Errorf("read file %s: %w", queuedFilepath, err)
	}

	var (
		currentTime     = time.Now()
		formattedTime   = currentTime.Format("2006-Jan-02")
		newFilenamePath = fmt.Sprintf("%sattempted_%s_%s", c.FilePath, formattedTime, queuedFilename)
	)

	if err := os.Rename(queuedFilepath, newFilenamePath); err != nil {
		return "", fmt.Errorf("rename from %s to %s: %w", queuedFilepath, newFilenamePath, err)
	}

	return content, nil
}
