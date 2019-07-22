// +build sqlite_json sqlite_json1 json1
// +build sqlite_fts5 fts5

package shenmeci

import (
	"database/sql"
	"encoding/json"
	"log"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func LoadDB() {
	config := GlobalConfig
	dbpath := filepath.Join(filepath.Dir(config.CedictPath), "shenmeci.sqlite")
	var err error
	db, err = sql.Open("sqlite3", dbpath)
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
		"PRAGMA quick_check",
		"SELECT json(1)",
	}
	for _, sql := range sqls {
		rows, err := db.Query(sql)
		if err != nil {
			log.Fatalf("%s: %s\n", sql, err)
		}
		defer rows.Close()
		for rows.Next() {
			var s string
			if err := rows.Scan(&s); err != nil {
				log.Fatal(err)
			}
			log.Printf("%s: %s", sql, s)
		}
		if err := rows.Err(); err != nil {
			log.Fatalf("%s: %s\n", sql, err)
		}
	}
	var hasTable bool
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='dict5'").Scan(&hasTable)
	if err != nil {
		log.Fatal(err)
	}
	if hasTable {
		log.Println("found FTS5 index")
	} else {
		log.Println("creating FTS5 index...")
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

func CloseDB() error {
	return db.Close()
}

func populateDB() error {
	start := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Do not call Rollback because journal_mode = OFF
	//defer tx.Rollback()

	// FIXME: set remove_diacritics=2 when ready to make database require
	// SQLite >= 3.27.0. See https://sqlite.org/releaselog/3_27_0.html and
	// https://sqlite.org/fts5.html#unicode61_tokenizer.
	// Previous SQLite versions throw an error when trying to operate on a
	// table with remove_diacritics=2.
	_, err = tx.Exec("CREATE VIRTUAL TABLE dict5 USING fts5(key, entry, tokenize = 'porter unicode61 remove_diacritics 1')")
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO dict5(key, entry) VALUES(?, ?)")
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
	rows, err := db.Query("SELECT key FROM dict5 WHERE entry MATCH ? ORDER BY rank", term)
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
