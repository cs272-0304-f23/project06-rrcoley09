package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"
)

func handleIndex(w http.ResponseWriter, r *http.Request, hits []Hits, imageUrl bool) {
	var t *template.Template
	var err error

	if len(hits) > 0 {
		t, err = template.ParseFiles("result.html")
	} else if len(hits) == 0 {
		t, err = template.ParseFiles("no_result.html")
	}

	if imageUrl {
		t, err = template.ParseFiles("image_result.html")
	}

	if err != nil {
		log.Fatalln("ParseFiles: ", err)
	}

	err = t.Execute(w, hits)
	if err != nil {
		log.Fatalln("Execute: ", err)
	}
}

func serve_web() {
	// Serve webpage with search form
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Defines handler for /search
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Get Search term and wildcard
		term := r.FormValue("term")
		wildcard := r.FormValue("wildcard")
		image := r.FormValue("image")
		var imageUrl bool

		sliceTerms := strings.Split(term, " ")

		// Initilaize slice of search results
		var hits []Hits

		// Check if the wildcard parameter is provided for the search.
		if len(wildcard) > 0 {
			// Perform a tf-idf search with wildcard functionality
			hits = search(sliceTerms, true, false)
		} else {
			if len(image) > 0 {
				imageUrl = true
				hits = search_image(term)
			} else {
				// Perform a tf-idf search without wildcard functionality
				hits = search(sliceTerms, false, false)
			}
		}

		handleIndex(w, r, hits, imageUrl)

	})
	// Start web server on port 8080
	// End connection if error (such as port in use)
	http.ListenAndServe(":8080", nil)
}
