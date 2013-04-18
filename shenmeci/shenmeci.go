package main

import (
	"log"
	"net/http"
	"os"
)

var (
	cedictPath = os.Getenv("CEDICT")
	cedict     *CEDICT
	err        error
)

func main() {
	cedict, err = loadCEDICT(cedictPath)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/segment", segmentHandler)
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
