package main

import (
	"database/sql"
	"log"
)

var (
	cedict *CEDICT
	db     *sql.DB
)

func main() {
	loadConfig()
	validateConfig()
	var err error
	cedict, err = loadCEDICT(config.CedictPath)
	if err != nil {
		log.Fatal(err)
	}
	db = createDB(cedict.Dict)
	defer db.Close()
	serve(config.Http.Host, config.Http.Port)
}
