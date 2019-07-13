package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func loadDB() {
	var err error
	db, err = sql.Open("sqlite3", filepath.Join(filepath.Dir(config.CedictPath), "shenmeci.sqlite"))
	if err != nil {
		log.Fatal(err)
	}

	var sqliteVersion string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SQLite version:", sqliteVersion)

	sqls := []string{
		"PRAGMA page_size = 4096",
		"PRAGMA synchronous = OFF",
		"PRAGMA journal_mode = OFF",
		"PRAGMA temp_store = MEMORY",
		"PRAGMA cache_size = -20480",
	}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatalf("%q: %s\n", err, sql)
		}
	}
	var hasTable bool
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='dict'").Scan(&hasTable)
	if err != nil {
		log.Fatal(err)
	}
	if hasTable {
		log.Println("found FTS index")
	} else {
		log.Println("creating FTS index...")
		err = populateDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS query(json)")
	if err != nil {
		log.Fatal(err)
	}
}

func populateDB() error {
	start := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Do not call Rollback because journal_mode = OFF
	//defer tx.Rollback()

	_, err = tx.Exec("CREATE VIRTUAL TABLE dict USING fts4(key, entry)")
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO dict(key, entry) VALUES(?, ?)")
	if err != nil {
		return err
	}
	dict := cedict.Dict
	for key, entry := range dict {
		for _, definition := range entry.definitions {
			stmt.Exec(key, definition)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	duration := time.Since(start).Truncate(time.Millisecond)
	log.Printf("indexed %d entries in %v\n", len(dict), duration)
	return nil
}

func searchDB(term string) (results []string) {
	rows, err := db.Query("SELECT key FROM dict WHERE entry MATCH ?", term)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		rows.Scan(&key)
		results = append(results, key)
	}
	return results
}

func insertQueryRecord(qr QueryRecord) error {
	b, err := json.Marshal(qr)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO query VALUES(json(?))", b)
	return err
}
