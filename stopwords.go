package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func stopWords() map[string]bool {
	// read form json file into slice of bytes
	fBytes, err := os.ReadFile("./stopwords-en.json")

	// Error handling for reading from json file
	if err != nil {
		fmt.Printf("error reading from stopwords file: %v\n", err.Error())
	}

	// Convert byte slice into string slice
	var words []string
	err = json.Unmarshal(fBytes, &words)

	// Error handling for conversion into string slice
	if err != nil {
		fmt.Printf("error reading from json: %v\n", err.Error())
	}

	// Make map for stop word for O(1) search
	stopWordsMap := make(map[string]bool)

	// Add word to map
	for _, word := range words {
		stopWordsMap[word] = true
	}

	// Return stopword map
	return stopWordsMap
}
