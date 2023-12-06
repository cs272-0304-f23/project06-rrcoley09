package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func serve_web() {
	// Serve webpage with search form
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Defines handler for /search
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Get Search term and wildcard
		term := r.FormValue("term")
		wildcard := r.FormValue("wildcard")

		sliceTerms := strings.Split(term, " ")

		start := time.Now()
		// Initilaize slice of search results
		var hits []Hits

		// Check if the wildcard parameter is provided for the search.
		if len(wildcard) > 0 {
			// Perform a tf-idf search with wildcard functionality
			hits = search(sliceTerms, true)
		} else {
			// Perform a tf-idf search without wildcard functionality
			hits = search(sliceTerms, false)
		}

		fmt.Fprintf(w, `<html>`)

		finish := time.Since(start).Seconds()
		// Print number of search results
		fmt.Fprintf(w, `<p> Legion | Searching for %v </p><br>`, term)
		if len(hits) != 0 {
			fmt.Fprintf(w, `<p> About %v results (%v seconds)</p>`, len(hits), finish)
		} else {
			fmt.Fprintf(w, `<p> No results found for search term....%v</p>`, term)
		}

		// Loop over map printing title with embedded href
		for _, hit := range hits {
			fmt.Fprintf(w, `<a href = %v> %v </a><br>`, hit.url, hit.title)
		}

		fmt.Fprintf(w, `</html>`)

	})
	// Start web server on port 8080
	// End connection if error (such as port in use)
	http.ListenAndServe(":8080", nil)
}
