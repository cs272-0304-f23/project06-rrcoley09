package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func createTables(pastDB bool) *sql.DB {
	// Open database
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if !pastDB {
		// Enable foreign key constraints
		_, err = db.Exec(`PRAGMA foreign_keys = ON;`)
		if err != nil {
			fmt.Printf("Error enabling foreign key constraints: %v\n", err)
		}

		// Create the "terms" table
		_, err = db.Exec(`DROP TABLE IF EXISTS terms;
		CREATE TABLE terms (
			term_id INTEGER PRIMARY KEY,
			term TEXT UNIQUE NOT NULL
		);`)
		if err != nil {
			fmt.Printf("Error creating terms table: %v\n", err)
		}

		// Create the "urls" table
		_, err = db.Exec(`DROP TABLE IF EXISTS urls;
		CREATE TABLE urls (
			url_id INTEGER PRIMARY KEY,
			name TEXT UNIQUE,
			title TEXT,
			word_count INTEGER
		);`)
		if err != nil {
			fmt.Printf("Error creating terms table: %v\n", err)
		}
		_, err = db.Exec(`DROP TABLE IF EXISTS snippets;
		CREATE TABLE snippets (
			snippet_id INTEGER PRIMARY KEY,
			sentence TEXT UNIQUE
		);`)

		// Create the "hits" table with foreign key constraints
		_, err = db.Exec(`DROP TABLE IF EXISTS hits;
		CREATE TABLE hits (
			hit_id INTEGER PRIMARY KEY,
			term_id INTEGER,
			url_id INTEGER,
			term_count INTEGER,
			snippet_id TEXT,
			UNIQUE (term_id, url_id),
			FOREIGN KEY (term_id) REFERENCES terms(term_id),
			FOREIGN KEY (url_id) REFERENCES urls(url_id),
			FOREIGN KEY (snippet_id) REFERENCES snippets(snippet_id)
		);`)
		if err != nil {
			fmt.Printf("Error creating hits table: %v\n", err)
		}

		// Create the "bigrams" table
		_, err = db.Exec(`DROP TABLE IF EXISTS bigrams;
		CREATE TABLE bigrams (
			bigram_id INTEGER PRIMARY KEY,
			term_id_1 INTEGER,
			term_id_2 INTEGER,
			url_id INTEGER,
			term_count INTEGER,
			snippet_id TEXT,
			UNIQUE (term_id_1, term_id_2, url_id),
			FOREIGN KEY (url_id) REFERENCES urls(url_id),
			FOREIGN KEY (snippet_id) REFERENCES snippets(snippet_id)
		);`)
		if err != nil {
			fmt.Printf("Error creating urls table: %v\n", err)
		}

		// Create the "caption_terms" table
		_, err = db.Exec(`DROP TABLE IF EXISTS caption_terms;
		CREATE TABLE caption_terms (
			term_id INTEGER PRIMARY KEY,
			term TEXT UNIQUE NOT NULL
		);`)
		if err != nil {
			fmt.Printf("Error creating caption_terms table: %v\n", err)
		}

		_, err = db.Exec(`DROP TABLE IF EXISTS hits_images;
		CREATE TABLE hits_images (
			image_id INTEGER PRIMARY KEY,
			term_id INTEGER,
			url_id INTEGER,
			image_url TEXT,
			term_count INTEGER,
			FOREIGN KEY (term_id) REFERENCES caption_terms(term_id),
			FOREIGN KEY (url_id) REFERENCES urls(url_id)
		);`)

	}
	return db
}
