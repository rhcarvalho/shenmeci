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
	db = createDB(cedict.Dict)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	serve(config.Http.Host, config.Http.Port)
}
