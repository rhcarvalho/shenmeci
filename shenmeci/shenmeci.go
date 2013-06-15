package main

import (
	"database/sql"
	"log"
	"os"
)

var (
	cedict *CEDICT
	db     *sql.DB
)

func main() {
	var err error
	cedictPath := os.Getenv("CEDICT")
	cedict, err = loadCEDICT(cedictPath)
	db = createDB(cedict.Dict)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	serve("127.0.0.1", os.Getenv("PORT"))
}
