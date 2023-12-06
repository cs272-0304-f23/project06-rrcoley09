package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
)

func insertSnippet(tx *sql.Tx, word string, word2 string, res ExtractResult) int {
	var snippetID int
	var bigram bool
	if word2 != "" {
		bigram = true
		word = word + " " + word2
	}

	for _, sentence := range res.sentences {
		if strings.Contains(sentence, word) {

			result, err := tx.Exec(`INSERT OR IGNORE INTO snippets (sentence) VALUES(?);`, sentence)

			if err != nil {
				fmt.Printf("error adding to snippets: %v\n", err)
				break
			}

			rowsAffected, _ := result.RowsAffected()
			if rowsAffected > 0 {
				lastInsertID, _ := result.LastInsertId()
				snippetID = int(lastInsertID)
			}

			if snippetID == 0 && bigram {
				err = tx.QueryRow(`SELECT snippet_id from snippets where sentence like ? limit 1;`, "%"+word+"%").Scan(&snippetID)
				if err != nil {
					if err != sql.ErrNoRows {
						fmt.Println("could not find snippet in table")
					}
				}

				//fmt.Printf("word: %v sentence: %v word: %v id: %v\n", word, sentence, word, snippetID)
			}

			break
		}
	}
	return snippetID
}

func add_to_table(res ExtractResult, db *sql.DB) {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer tx.Commit()

	// Insert url, title, and url word count into urls table
	result, err := tx.Exec(`INSERT OR IGNORE INTO urls (name, title, word_count) VALUES(? , ?, ?);`, res.url, res.webPageTitle, len(res.words))
	if err != nil {
		fmt.Printf("error inserting into urls table: %v\n", err)
	}
	// Retrieve the last inserted url_id
	urlID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("could not find URL in urls table")
	}

	mapWordIDs := make(map[string]int64)

	// Range over words and stem them
	for _, word := range res.words {
		_, err := tx.Exec(`INSERT OR IGNORE INTO terms (term)
									VALUES(?);`, word)
		if err != nil {
			fmt.Printf("error inserting term into table: %v\n", err)
		}

		// Query the term_id from the terms table
		var termID int
		err = tx.QueryRow(`SELECT term_id from terms where term = ?;`, word).Scan(&termID)
		if err != nil {
			fmt.Println("could not find term in terms table")
		}
		mapWordIDs[word] = int64(termID)

		snippetID := insertSnippet(tx, word, "", res)

		var hitID int
		err = tx.QueryRow(`SELECT hit_id from hits where term_id = ? and url_id = ?;`, termID, urlID).Scan(&hitID)
		if err != nil {
			if err == sql.ErrNoRows {
				//fmt.Printf("termID: %d, urlID: %d, snippetID: %d\n", termID, urlID, snippetID)
				if snippetID != 0 {
					_, err = tx.Exec(`INSERT OR IGNORE INTO hits (term_id, url_id, term_count, snippet_id) VALUES(?, ?, 1, ?);`, termID, urlID, snippetID)
					//fmt.Println(snippetID)
					if err != nil {
						fmt.Printf("error inserting new freq count in hits table with snippet: %v\n", err)
					}
				} else {
					_, err = tx.Exec(`INSERT OR IGNORE INTO hits (term_id, url_id, term_count) VALUES(?, ?, 1);`, termID, urlID)

					if err != nil {
						fmt.Printf("error inserting new freq count in hits table without snippet: %v\n", err)
					}
				}
			} else {
				fmt.Printf("error getting hit_id from hits table: %v\n", err)
			}
		} else {
			_, err = tx.Exec(`UPDATE hits SET term_count = term_count + 1 WHERE hit_id = ?;`, hitID)
			if err != nil {
				fmt.Printf("error updating hits table: %v\n", err)
			}
		}

	}

	// Query for the url_id of the inserted URL
	for index, word := range res.words {
		var urlID int64
		err = tx.QueryRow(`SELECT url_id from urls where name = ?;`, res.url).Scan(&urlID)

		if err != nil {
			fmt.Println("could not find url in urls table")
		}

		termID := mapWordIDs[word]

		if index+1 < len(res.words) {
			word2 := res.words[index+1]
			termID2 := mapWordIDs[word2]

			snippetID := insertSnippet(tx, word, word2, res)
			//fmt.Printf("bigram: %v\n", snippetID)

			var bigramID int
			err = tx.QueryRow(`SELECT bigram_id from bigrams where term_id_1 = ? and term_id_2 = ? and url_id = ?;`, termID, termID2, urlID).Scan(&bigramID)
			if err != nil {
				if err == sql.ErrNoRows {
					if snippetID != 0 {
						// If no rows are found, insert a new row into the "hits" table
						_, err = tx.Exec(`INSERT INTO bigrams (term_id_1, term_id_2, url_id, term_count, snippet_id) VALUES(?, ?, ?, 1, ?);`, termID, termID2, urlID, snippetID)

						if err != nil {
							fmt.Println("error inserting new freq count in bigrams table with snippet")
						}
					} else {
						_, err = tx.Exec(`INSERT INTO bigrams (term_id_1, term_id_2, url_id, term_count) VALUES(?, ?, ?, 1);`, termID, termID2, urlID)

						if err != nil {
							fmt.Println("error inserting new freq count in bigrams table without snippet")
						}
					}
				} else {
					fmt.Printf("error getting bigram_id from bigrams table: %v\n", err)
				}
			} else {
				// If a row is found, update the term_count in the "hits" table
				_, err = tx.Exec(`UPDATE bigrams SET term_count = term_count + 1 WHERE bigram_id = ?;`, bigramID)
			}
		}
	}

	for src, alt := range res.images {
		for _, word := range alt {
			_, err := tx.Exec(`INSERT OR IGNORE INTO caption_terms (term)
									VALUES(?);`, word)
			if err != nil {
				fmt.Printf("error inserting term into table - caption_terms: %v\n", err)
			}

			var termID int
			err = tx.QueryRow(`SELECT term_id from caption_terms where term = ?;`, word).Scan(&termID)
			if err != nil {
				fmt.Println("could not find term in caterms table")
			}

			var imageID int
			err = tx.QueryRow(`SELECT image_id from hits_images where term_id = ? and image_url = ? and url_id = ?;`, termID, src, urlID).Scan(&imageID)
			if err != nil {
				if err == sql.ErrNoRows {
					_, err = tx.Exec(`INSERT OR IGNORE INTO hits_images (term_id, url_id, image_url, term_count) VALUES(?, ?, ?, 1);`, termID, urlID, src)
					if err != nil {
						fmt.Printf("error inserting new freq count in hits_images table: %v\n", err)
					}
				} else {
					fmt.Printf("error getting image_id from hits_images table: %v\n", err)
				}
			} else {
				_, err = tx.Exec(`UPDATE hits_images SET term_count = term_count + 1 WHERE image_id = ? and url_id = ?;`, imageID, urlID)
				if err != nil {
					fmt.Printf("error updating hits_images table: %v\n", err)
				}
			}
		}
	}
}

func populate_tables(exOutCh chan ExtractResult, wg *sync.WaitGroup, n *int32, pastDB bool) {
	// Create database tables if they don't exist
	db := createTables(pastDB)

	// Loop over results from extract for each url
	for extractResult := range exOutCh {
		// Add extracted result to the database
		add_to_table(extractResult, db)

		// Decrements wait group counter by 1
		atomic.AddInt32(n, -1)
		wg.Done()
	}
	db.Close()
}
