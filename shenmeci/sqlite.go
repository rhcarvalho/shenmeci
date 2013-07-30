package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path"
)

var db *sql.DB

func loadDB() {
	db, err := sql.Open("sqlite3", path.Join(path.Dir(config.CedictPath), "shenmeci.sqlite"))
	if err != nil {
		log.Fatal(err)
	}

	var sqliteVersion string
	err = db.QueryRow("select sqlite_version()").Scan(&sqliteVersion)
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
		populateDB()
	}
}

func populateDB() {
	_, err := db.Exec("CREATE VIRTUAL TABLE dict USING fts4(key, entry)")
	if err != nil {
		log.Fatal(err)
	}
	insertStmt, err := db.Prepare("INSERT INTO dict(key, entry) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStmt.Close()
	dict := cedict.Dict
	for key, entry := range dict {
		for _, definition := range entry.definitions {
			insertStmt.Exec(key, definition)
		}
	}
	log.Printf("indexed %d entries\n", len(dict))
}

func searchDB(db *sql.DB, term string) (results []string) {
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
