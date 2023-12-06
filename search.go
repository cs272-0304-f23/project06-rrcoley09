package main

import (
	"fmt"

	"github.com/kljensen/snowball"
)

func search(terms []string, wildCard bool) []Hits {
	// Stem the search term
	stemmedSearch := []string{}

	for _, searchTerm := range terms {
		stemmed, err := snowball.Stem(searchTerm, "english", true)
		if err != nil {
			fmt.Printf("error steming search term in search_tf-idf: %v\n", err)
		}
		stemmedSearch = append(stemmedSearch, stemmed)
	}

	// If stemming is successful, perform tf-idf search with the stemmmed term and wildcard flag
	return tfIdf(stemmedSearch, wildCard)
}
