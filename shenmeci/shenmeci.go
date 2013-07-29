package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
)

var (
	cedict *CEDICT
	db     *sql.DB
)

func main() {
	loadConfig()
	validateConfig()

	// Test whether we can listen on the provided Host and Port.
	// If the Host:Port is already in use, we can exit before wasting more resources.
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Http.Host, config.Http.Port))
	if err != nil {
		log.Fatal(err)
	}
	ln.Close()

	cedict, err = loadCEDICT(config.CedictPath)
	if err != nil {
		log.Fatal(err)
	}
	db = createDB(cedict.Dict)
	defer db.Close()
	serve(config.Http.Host, config.Http.Port)
}
