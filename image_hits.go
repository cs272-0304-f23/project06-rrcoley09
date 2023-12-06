package main

import (
	"fmt"

	"github.com/kljensen/snowball"
)

func search_image(term string) []Hits {
	stemmed, err := snowball.Stem(term, "english", true)
	if err != nil {
		fmt.Printf("error steming search term in search_image: %v\n", err)
	}

	stemmedSlice := []string{stemmed}

	return tfIdf(stemmedSlice, false, true)
}


