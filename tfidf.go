package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
)

// Instantiate struct holding url and tf-idf rank
type Hits struct {
	Title    string
	Image    string
	Sentence string
	Url      string
	Score    float64
}

// Create slice of struct values
type ByRank []Hits

// Implement functions of the sort interface
func (hits ByRank) Len() int {
	return len(hits)
}

func (hits ByRank) Less(i, j int) bool {
	return hits[i].Score > hits[j].Score
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

func get_images(term string, numOfDocs int, db *sql.DB) []Hits {
	var docsWithTerm int
	sliceTFIDF := []Hits{}
	// Get term id for input each term
	var termID int
	err := db.QueryRow(`SELECT term_id 
					FROM caption_terms 
					WHERE term = ?;`, term).Scan(&termID)
	if err != nil {
		fmt.Printf("error getting term id for search term: %v\n", err)
	}

	// Count number of docs with given term using query
	err = db.QueryRow(`SELECT count(*) 
					FROM hits_images 
					WHERE term_id = ?;`, termID).Scan(&docsWithTerm)
	if err != nil {
		fmt.Printf("error getting count of docs with term: %v\n", err)
	}

	// Join hits and urls on url_id, get resulting attributes for search term
	rows, err := db.Query(`SELECT urls.name, urls.title, hits_images.image_url, hits_images.term_count, urls.word_count
						FROM urls 
						INNER JOIN hits_images 
						WHERE urls.url_id = hits_images.url_id 
						AND hits_images.term_id = ?;`, termID)
	if err != nil {
		fmt.Printf("error getting results from join of hits_images and urls: %v\n", err)
	}

	for rows.Next() {
		// Initialize variables and scan values for attributes extratced from table
		var urlName, urlTitle, imageUrl string
		var termCount, termsInDoc int
		err = rows.Scan(&urlName, &urlTitle, &imageUrl, &termCount, &termsInDoc)
		if err != nil {
			fmt.Println("images: error getting values for results from hits join urls")
		}

		// Calculate tf-idf score for testing
		tfIdf := calcTfIdf(termCount, termsInDoc, docsWithTerm, numOfDocs)

		// Add struct to slice
		rankURL := Hits{urlTitle, imageUrl, "", urlName, tfIdf}

		// Append struct to slice
		sliceTFIDF = append(sliceTFIDF, rankURL)
	}

	return sliceTFIDF
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
		rows, err = db.Query(`SELECT urls.name, urls.title, hits.term_count, urls.word_count, hits.snippet_id
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
			if err != sql.ErrNoRows {
				fmt.Printf("error getting term id for search term w: %v\n", err)
			}
		}

		// Count number of docs with given term using query
		err = db.QueryRow(`SELECT count(*) 
						FROM hits 
						WHERE term_id = ?;`, termID).Scan(&docsWithTerm)
		if err != nil {
			fmt.Printf("error getting count of docs with term: %v\n", err)
		}

		// Join hits and urls on url_id, get resulting attributes for search term
		rows, err = db.Query(`SELECT urls.name, urls.title, hits.term_count, urls.word_count, hits.snippet_id
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
	rows, err = db.Query(`SELECT urls.name, urls.title, bigrams.term_count, urls.word_count, bigrams.snippet_id
						FROM urls 
						INNER JOIN bigrams 
						WHERE urls.url_id = bigrams.url_id AND (bigrams.term_id_1 = ? AND bigrams.term_id_2 = ?);`, termID, termID2)
	if err != nil {
		fmt.Printf("bigrams: error getting results from join of hits and urls: %v\n", err)
	}

	return rows, docsWithTerm
}

func tfIdf(terms []string, wildCard bool, imageSearch bool) []Hits {
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

	if imageSearch {
		term := terms[0]
		sliceTFIDF = get_images(term, numOfDocs, db)
	} else {
		if len(terms) == 1 {
			term := terms[0]
			rows, docsWithTerm = get_hits(term, wildCard, db)
		} else {
			rows, docsWithTerm = get_bigrams(terms, db)
		}
	}

	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			// Initialize variables and scan values for attributes extratced from table
			var urlName, urlTitle string
			var termCount, termsInDoc int
			var snippetID sql.NullInt64
			err = rows.Scan(&urlName, &urlTitle, &termCount, &termsInDoc, &snippetID)
			if err != nil {
				fmt.Printf("error scanning rows from table join: %v\n", err)
			}

			var sentence string
			if snippetID.Valid {
				err = db.QueryRow(`SELECT sentence 
                    FROM snippets
                    WHERE snippet_id = ?;`, snippetID.Int64).Scan(&sentence)
				if err != nil {
					fmt.Println("error querying for snippet sentence")
				}
			}
			// Calculate tf-idf score for testing
			tfIdf := calcTfIdf(termCount, termsInDoc, docsWithTerm, numOfDocs)

			// Add struct to slice
			rankURL := Hits{urlTitle, "", sentence, urlName, tfIdf}

			// Append struct to slice
			sliceTFIDF = append(sliceTFIDF, rankURL)
		}
	}
	// Sort struct, order descending
	sort.Sort(ByRank(sliceTFIDF))

	// Return slice of structs
	return sliceTFIDF
}
