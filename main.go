package main

import (
	"fmt"
	"log"
	"os"
	"time"

	// other imports
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	fmt.Println("Go-Twitter Bot v0.01")
	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, err := getClient(&creds)
	if err != nil {
		log.Panicf("Error getting Twitter Client: %s", err.Error())
		return
	}

	person := "@LudwigAhgren"
	tweetID := int64(1342270499825414145)
	tweetContent := fmt.Sprintf("%s You're so tall! https://www.amazon.com/hz/wishlist/ls/3CJCHZ9DGN05D?ref_=wl_share check it out my friend!%d", person, ticks)

	ticks := 0
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			// for now ignore tweet and status
			_, _, err := client.Statuses.Update(tweetContent, &twitter.StatusUpdateParams{
			InReplyToStatusID: tweetID,
			})

			if err != nil {
				panic(err)
			}

			fmt.Println("Tweet sent! ", ticks)
		}

		ticks++
	}
}


// Credentials stores all of our access/consumer tokens
// and secret keys needed for authentication against
// the twitter REST API.
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// getClient is a helper function that will return a twitter client
// that we can subsequently use to send tweets, or to stream new tweets
// this will take in a pointer to a Credential struct which will contain
// everything needed to authenticate and return a pointer to a twitter Client
// or an error
func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	_, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	return client, nil
}