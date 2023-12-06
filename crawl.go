package main

import (
	"sync"
	"sync/atomic"
)

func crawl(seedURL string, siteMap bool, pastDB bool) {
	// Create channels to pass information in goroutines
	dInCh := make(chan string, 10000)
	dOutCh := make(chan DownloadResult, 10000)
	exOutCh := make(chan ExtractResult, 10000)

	// Create wait group variable and atomic counter variable
	var n int32
	wg := sync.WaitGroup{}

	if !pastDB {
		// Input seed url into download channel to be crawled first
		dInCh <- seedURL

		// Add seed url to wait gropp and increment atomic counter
		wg.Add(1)
		atomic.AddInt32(&n, 1)

		// Create multiple threads existing outside main thread using goroutines
		// Create alternate control flow based on boolean value sitemap
		if siteMap {
			// If there is a site map, download and extract take in true boolean
			go download(dInCh, dOutCh, &wg, &n, true)
			go extract(dInCh, dOutCh, exOutCh, &wg, &n, true)
		} else {
			// If there is no site map, download and extract take in false boolean
			go download(dInCh, dOutCh, &wg, &n, false)
			go extract(dInCh, dOutCh, exOutCh, &wg, &n, false)
		}

		// Start the index update process in a separate goroutine
		go populate_tables(exOutCh, &wg, &n, pastDB)
	}

	// Block until wait group counter is zero
	// No more urls to crawl, extract, and add to index
	wg.Wait()

	// Check if the atomic counter is zero before closing channels
	if atomic.LoadInt32(&n) == 0 {
		close(dInCh)
		close(dOutCh)
		close(exOutCh)
	}

}
