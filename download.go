package main

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const defaultUserAgent = "Go-http-client/1.1"

type DownloadResult struct {
	url  string
	body []byte
}

func checkDisallow(url string) (bool, AgentRules) {
	record := get_robots(url)

	for agent, agentRules := range record.userAgent {
		// Replace "*" in the agent pattern with ".*" for regex matching
		agentPattern := strings.ReplaceAll(agent, "*", ".*")
		// Check if the default user agent matches the agent pattern
		checkAgent, _ := regexp.Match(agentPattern, []byte(defaultUserAgent))
		if checkAgent {
			for val := range agentRules.disallow {
				// Check if the url matches any disallow pattern
				match, _ := regexp.MatchString(val, url)
				if match {
					return true, agentRules
				}
			}
		}
	}
	return false, AgentRules{}
}

func download(inC chan string, dOutCh chan DownloadResult, wg *sync.WaitGroup, n *int32, siteMap bool) {
	// Loop over URLs in inC
	for url := range inC {
		// Check if the url is disallowed and get the agent rules
		match, record := checkDisallow(url)

		// Sleep for crawl delay if specified in agent rules
		if record.crawlDelay > 0 {
			time.Sleep(time.Duration(record.crawlDelay) * time.Second)
		}

		// Proceed with the download if not disallowed
		if !match {
			resp, err := http.Get(url)
			if err == nil {
				// Ensure the response of HTTP request is closed
				defer resp.Body.Close()

				// Check the success status of the request
				if resp.StatusCode == 200 {
					// Reads from io.Reader and returns info as a slice of bytes and error
					body, err := io.ReadAll(resp.Body)
					if err == nil {
						// Return download result to dOutCh for use in extract
						dOutCh <- DownloadResult{url, body}
					}
				}
			}
		} else {
			// Decrement the counter and signal the wait group
			atomic.AddInt32(n, -1)
			wg.Done()
		}
	}
}
