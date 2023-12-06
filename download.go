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
		agentPattern := strings.ReplaceAll(agent, "*", ".*")
		checkAgent, _ := regexp.Match(agentPattern, []byte(defaultUserAgent))
		if checkAgent {
			for val := range agentRules.disallow {
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
	// Loop over urls in inC
	for url := range inC {
		match, record := checkDisallow(url)

		if record.crawlDelay > 0 {
			time.Sleep(time.Duration(record.crawlDelay) * time.Second)
		}

		if !match {
			resp, err := http.Get(url)
			if err == nil {
				// Ensure resp of http request is closed
				defer resp.Body.Close()

				// Check success status of request
				if resp.StatusCode == 200 {
					// Reads from io.Reader and returns info as slice of bytes and error
					body, err := io.ReadAll(resp.Body)
					if err == nil {
						// Return download result to dOutCh for use in extract
						dOutCh <- DownloadResult{url, body}
					}
				}
			}
		} else {
			atomic.AddInt32(n, -1)
			wg.Done()
		}
	}
}
