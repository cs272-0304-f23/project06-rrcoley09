package main

import (
	"bufio"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Records struct {
	userAgent map[string]AgentRules
}

type AgentRules struct {
	disallow   map[string]struct{}
	crawlDelay int
}

func get_robots(seedURL string) Records {
	// Parse seed url
	parsedURL, err := url.Parse(seedURL)

	var record Records
	record.userAgent = make(map[string]AgentRules)

	if err == nil {
		// Set the path to the robots.txt file
		parsedURL.Path = "robots.txt"
		resp, err := http.Get(parsedURL.String())
		if err == nil {
			defer resp.Body.Close()

			// Check if the response status code is 200 (OK)
			if resp.StatusCode == 200 {
				scanner := bufio.NewScanner(resp.Body)

				var currAgent string
				var agentRules AgentRules

				// Iterate over each line in the robots.txt file
				for scanner.Scan() {
					line := scanner.Text()

					// Check for "User-agent:" directive
					if strings.Contains(line, "User-agent: ") {
						if currAgent != "" {
							record.userAgent[currAgent] = agentRules
						}
						splitLine := strings.Split(line, "User-agent: ")
						currAgent = splitLine[1]
						agentRules = AgentRules{
							disallow: make(map[string]struct{}),
						}
					} else if strings.Contains(line, "Disallow: ") {
						// Check for "Disallow:" directive
						splitLine := strings.Split(line, "Disallow: ")

						// Add disallowed path to the map
						pattern := strings.ReplaceAll(splitLine[1], "*", ".*")
						agentRules.disallow[pattern] = struct{}{}
					} else if strings.Contains(line, "Crawl-delay: ") {
						// Check for "Crawl-delay:" directive
						splitLine := strings.Split(line, "Crawl-delay: ")

						// Convert crawl delay to integer and assign to the record
						agentRules.crawlDelay, _ = strconv.Atoi(splitLine[1])
					}
				}
				if currAgent != "" {
					record.userAgent[currAgent] = agentRules
				}
			}
		}
	}
	return record
}
