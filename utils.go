package main

import (
	"encoding/json"
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
// from ./credentials.json - the file should be created
// following format in credentials_example
func loadCredentials() {
	file, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalln("Could not load credentials.json file (See credentials_example.json):", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalln("Error parsing credentials.json, please confirm valid json format", err)
	}
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
