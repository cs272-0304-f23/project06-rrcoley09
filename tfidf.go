package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
)

// Instantiate struct holding url and tf-idf rank
type Hits struct {
	title string
	url   string
	score float64
}

// Create slice of struct values
type ByRank []Hits

// Implement functions of the sort interface
func (hits ByRank) Len() int {
	return len(hits)
}

func (hits ByRank) Less(i, j int) bool {
	return hits[i].score > hits[j].score
}

func (hits ByRank) Swap(i, j int) {
	hits[i], hits[j] = hits[j], hits[i]
}

func calcTfIdf(termCount, termsInDoc, docsWithTerm, numOfDocs int) float64 {
	// Tf - count of terms / total terms in doc
	tf := float64(termCount) / float64(termsInDoc)

	// Df - docs with term occurrence / total number of docs
	// Idf - inverse of Df
	df := float64(docsWithTerm) / float64(numOfDocs)
	idF := 1 / df

	// Doc relevancy to search term score
	tfIdf := tf * idF

	return tfIdf
}

func get_hits(term string, wildCard bool, db *sql.DB) (*sql.Rows, int) {
	var docsWithTerm int
	var rows *sql.Rows

	if wildCard {
		// Query for counting the number of documents with the given term using wildcard search
		err := db.QueryRow(`SELECT count(*) FROM hits
				INNER JOIN terms ON hits.term_id = terms.term_id
				WHERE terms.term LIKE ?;`, term+"%").Scan(&docsWithTerm)
		if err != nil {
			fmt.Printf("error getting count of docs with term: %v\n", err)
		}

		// Retrieve results from a join of hits, terms, and urls based on the wildcard search
		rows, err = db.Query(`SELECT urls.name, urls.title, hits.term_count, urls.word_count
						FROM hits
						INNER JOIN terms ON hits.term_id = terms.term_id
						INNER JOIN urls ON hits.url_id = urls.url_id 
						WHERE terms.term LIKE ?;`, term+"%")
		if err != nil {
			fmt.Printf("error getting results from join of hits and urls: %v\n", err)
		}
	} else {
		// Get term id for input each term
		var termID int
		err := db.QueryRow(`SELECT term_id 
						FROM terms 
						WHERE term = ?;`, term).Scan(&termID)
		if err != nil {
			fmt.Printf("error getting term id for search term: %v\n", err)
		}

		// Count number of docs with given term using query
		err = db.QueryRow(`SELECT count(*) 
						FROM hits 
						WHERE term_id = ?;`, termID).Scan(&docsWithTerm)
		if err != nil {
			fmt.Printf("error getting count of docs with term: %v\n", err)
		}

		// Join hits and urls on url_id, get resulting attributes for search term
		rows, err = db.Query(`SELECT urls.name, urls.title, hits.term_count, urls.word_count
							FROM urls 
							INNER JOIN hits 
							WHERE urls.url_id = hits.url_id 
							AND hits.term_id = ?;`, termID)
		if err != nil {
			fmt.Printf("error getting results from join of hits and urls: %v\n", err)
		}
	}

	return rows, docsWithTerm
}

func get_bigrams(terms []string, db *sql.DB) (*sql.Rows, int) {
	var docsWithTerm int
	var rows *sql.Rows

	// Get term id for input each term
	var termID int
	err := db.QueryRow(`SELECT term_id 
					FROM terms 
					WHERE term = ?;`, terms[0]).Scan(&termID)

	var termID2 int
	err = db.QueryRow(`SELECT term_id 
					FROM terms 
					WHERE term = ?;`, terms[1]).Scan(&termID2)

	// Count number of docs with given term using query
	err = db.QueryRow(`SELECT count(*) 
					FROM bigrams 
					WHERE term_id_1 = ? AND term_id_2 = ?;`, termID, termID2).Scan(&docsWithTerm)
	if err != nil {
		fmt.Printf("bigrams: error getting count of docs with term: %v\n", err)
	}

	// Join hits and urls on url_id, get resulting attributes for search term
	rows, err = db.Query(`SELECT urls.name, urls.title, bigrams.term_count, urls.word_count
						FROM urls 
						INNER JOIN bigrams 
						WHERE urls.url_id = bigrams.url_id AND (bigrams.term_id_1 = ? AND bigrams.term_id_2 = ?);`, termID, termID2)
	if err != nil {
		fmt.Printf("bigrams: error getting results from join of hits and urls: %v\n", err)
	}

	return rows, docsWithTerm
}

func tfIdf(terms []string, wildCard bool) []Hits {
	// Create slice for return value
	sliceTFIDF := []Hits{}

	// Open db and defer close until function exit
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(-1)
	}

	defer db.Close()

	// Query rows in urls table and find count of urls as number of docs
	var numOfDocs int
	err = db.QueryRow(`SELECT count(*) 
					FROM urls;`).Scan(&numOfDocs)
	if err != nil {
		fmt.Printf("error getting num of docs (url table): %v\n", err)
	}

	var docsWithTerm int
	var rows *sql.Rows

	if len(terms) == 1 {
		term := terms[0]
		get_hits(term, wildCard, db)
		rows, docsWithTerm = get_hits(term, wildCard, db)
	} else {
		rows, docsWithTerm = get_bigrams(terms, db)
	}

	for rows.Next() {
		// Initialize variables and scan values for attributes extratced from table
		var urlName, urlTitle string
		var termCount, termsInDoc int
		err = rows.Scan(&urlName, &urlTitle, &termCount, &termsInDoc)
		if err != nil {
			fmt.Println("bigrams: error getting values for results from hits join urls")
		}

		// Calculate tf-idf score for testing
		tfIdf := calcTfIdf(termCount, termsInDoc, docsWithTerm, numOfDocs)

		// Add struct to slice
		rankURL := Hits{urlTitle, urlName, tfIdf}

		// Append struct to slice
		sliceTFIDF = append(sliceTFIDF, rankURL)
	}
	// Sort struct, order descending
	sort.Sort(ByRank(sliceTFIDF))

	// Return slice of structs
	return sliceTFIDF
}
