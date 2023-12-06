package main

import (
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
)

type (
	// Sitemap represents information about individual sitemap
	Sitemap struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	}

	// SitemapIndex represents the sitemap index structure
	SitemapIndex struct {
		Sitemaps []Sitemap `xml:"sitemap"`
	}

	// Url represents information about a single url entry within a sitemap
	Url struct {
		Loc        string `xml:"loc"`
		LastMod    string `xml:"lastmod"`
		ChangeFreq string `xml:"changefreq"`
		Priority   string `xml:"priority"`
	}

	// Urlset represents the structure of a sitemap containing multiple urls
	Urlset struct {
		Urls []Url `xml:"url"`
	}
)

func extractSitemapUrl(body []byte, n *int32, wg *sync.WaitGroup, inC chan string) {
	// Parse the sitemap index
	var siteMapData SitemapIndex
	xml.Unmarshal(body, &siteMapData)

	// Iterate over each sitemap in the index
	for _, siteMap := range siteMapData.Sitemaps {
		var urlSetData Urlset

		// Download the body of the individual sitemap
		body, err := getSitemapBody(siteMap.Loc)
		if err != nil {
			fmt.Printf("Error downloading individual sitemap: %v\n", err)
		}

		// Unescape html entities in the sitemap body
		body = []byte(html.UnescapeString(string(body)))

		// Unmarshal the sitemap body to extract url information
		err = xml.Unmarshal(body, &urlSetData)
		if err != nil {
			fmt.Printf("Error unmarshalling Sitemap: %v\n", err)
		}

		// Iterate over each url in the sitemap and send it to the input channel
		for _, url := range urlSetData.Urls {
			wg.Add(1)
			atomic.AddInt32(n, 1)
			inC <- url.Loc
		}
	}
}

func getSitemapBody(url string) ([]byte, error) {
	// Issues GET to specified url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read the body of the http response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
