package main

import (
	"flag"
	"fmt"
	"github.com/darkhelmet/twitterstream"
	"github.com/njern/gogmail"
	"log"
	"runtime"
	"sort"
	"time"
)

const (
	RECONNECT_TIME     = time.Duration(30 * time.Second)
	LINK_SEND_INTERVAL = time.Duration(24 * time.Hour)
)

var (
	topics            = flag.String("topics", "", "The topics Hutch should be tracking.")
	numberOfLinks     = flag.Int("num_links", 5, "The number of unique links per topic to track (daily)")
	trackedLinks      = make(map[string]int)
	lastLinksSentTime = time.Now()
	gmail             *gogmail.GMail
)

func parseFlags() {
	flag.Parse()
	if *topics == "" {
		log.Fatalln("Please specify --topics \"some topics here\"(according to https://dev.twitter.com/docs/streaming-apis/parameters#track)")
	}
}

func handleIncomingTweet(tweet *twitterstream.Tweet) {
	if len(tweet.Entities.Urls) == 0 {
		//log.Printf("Skipping tweet: '%s', written by %s - It does not contain any URL's!\n", tweet.Text, tweet.User.ScreenName)
		return
	}

	for _, url := range tweet.Entities.Urls {
		finalURL, err := resolveFinalURL(url.Url)
		if err != nil {
			log.Printf("Error resolving URL %s: Error was: %s\n", url, err)
			return
		}

		// Store the number of times we've seen this particular, unique URL
		_, ok := trackedLinks[finalURL]
		if !ok {
			trackedLinks[finalURL] = 1
		} else {
			trackedLinks[finalURL] += 1
		}
	}

	log.Printf("Time since lastLinksSent: %s", time.Since(lastLinksSentTime))
	// Send the most popular links and clear the list every twenty-four hours.
	if time.Since(lastLinksSentTime).Seconds() > LINK_SEND_INTERVAL.Seconds() {
		// Refresh the lastLinksSentTime
		lastLinksSentTime = time.Now()

		// Extract the list of unique links
		var linkScores []LinkScore
		for link, score := range trackedLinks {
			linkScores = append(linkScores, LinkScore{link, score})
		}

		// Sort in descending order
		sort.Sort(sort.Reverse(ByScore(linkScores)))

		// Send the top links (in order) via e-mail (TODO: HTML template)
		list := fmt.Sprintf("Hi there %s, here are your links for the last 24 hours:\n\n\n", gmail.Username)
		for i := 0; i < min(*numberOfLinks, len(linkScores)); i++ {
			list += fmt.Sprintf("\t* %s - %d mentions\n", linkScores[i].link, linkScores[i].score)
		}
		list += "\nKind regards,\n\nHutch"
		subject := fmt.Sprintf("Your daily tracked links from Twitter")

		log.Println(list)
		err := gmail.SendMail([]string{gmail.Username}, subject, list, false)
		if err != nil {
			log.Fatalf("Something went horribly wrong sending your daily e-mail! Error was: %s\n", err)
		}

		log.Fatalf("Sent mail to %s!\n", gmail.Username)

		// Empty the list and begin the dance all over again.
		trackedLinks = make(map[string]int)
	}
}

func init() {
	parseFlags() // Parse flags
}

func main() {
	// Load credentials
	consumerKey, consumerSecret, accessToken, accessSecret, mailAddress, gmailPassword := loadCredentials()
	// Initialize Twitter streaming client.
	twitterStream := twitterstream.NewClient(consumerKey, consumerSecret, accessToken, accessSecret)
	// Initialize e-mail "client"
	gmail = gogmail.GmailConnection(mailAddress, gmailPassword)

	for {
		stream, err := twitterStream.Track(*topics)
		if err != nil {
			log.Printf("Connecting to the streaming API failed, reconnecting in %s...\n", RECONNECT_TIME)
			time.Sleep(RECONNECT_TIME)
			continue
		}

		for {
			tweet, err := stream.Next()
			if err != nil {
				log.Println("Connection died with error:", err)
				log.Printf("Reconnecting in %s...\n", RECONNECT_TIME)
				break
			}
			// Handle the Tweet that just came in
			handleIncomingTweet(tweet)
		}

		// We sleep a while before reconnecting to keep Twitter happy.
		time.Sleep(RECONNECT_TIME)
	}
}
