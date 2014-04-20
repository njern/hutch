package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// loadCredentials() loads the user's credentials
// from ./CREDENTIALS - the file should be created
// as in CREDENTIALS_EXAMPLE
func loadCredentials() (consumerKey, consumerSecret, accessToken, accessSecret, gmailUser, gmailPassword string) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	if err != nil {
		log.Fatalln("Could not load CREDENTIALS file (See CREDENTIALS_EXAMPLE for an example):", err)
	}

	lines := strings.Split(string(credentials), "\n")
	if len(lines) < 6 {
		log.Fatalln("Your CREDENTIALS file should have at least six line (look at CREDENTIALS_EXAMPLE)")
	}

	consumerKey = lines[0]
	consumerSecret = lines[1]
	accessToken = lines[2]
	accessSecret = lines[3]
	gmailUser = lines[4]
	gmailPassword = lines[5]

	return
}

// resolveFinalURL resolves a given URL. Since all Twitter
// links are shortened at least once, we need to resolve
// the URL in order for it to be comparable.
func resolveFinalURL(initialURL string) (string, error) {
	resp, err := http.Get(initialURL)
	if err != nil {
		return "", err
	}

	// The Request in the Response is the last
	// URL the client tried to access.
	finalURL := resp.Request.URL.String()
	resp.Body.Close()

	// Cut out any GET variables in the URL by removing everything after a "?"
	parts := strings.Split(finalURL, "?")

	return parts[0], nil
}
