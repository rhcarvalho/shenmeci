package main

import (
	"log"
	"os"
)

var cedict *CEDICT

func main() {
	var err error
	cedictPath := os.Getenv("CEDICT")
	cedict, err = loadCEDICT(cedictPath)
	if err != nil {
		log.Fatal(err)
	}
	serve("127.0.0.1", os.Getenv("PORT"))
}
