package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	cedictPath = os.Getenv("CEDICT")
	cedict     *CEDICT
	err        error
)

func main() {
	if cedictPath == "" {
		// if the environment variable is not set, try a default path
		cedictPath = "dict/cedict_1_0_ts_utf-8_mdbg.txt.gz"
		if _, err := os.Stat(cedictPath); err != nil {
			if os.IsNotExist(err) {
				// try to download
				cmd := exec.Command("./download_dict.sh")
				err = cmd.Run()
				if err != nil {
					// file does not exist and could not be downloaded
					log.Fatal("Missing environment variable CEDICT.")
				}
			} else {
				// other error
				log.Fatal(err)
			}
		}
	}
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
