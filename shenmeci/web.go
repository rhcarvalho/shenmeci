package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func segmentHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	results := keysToResults(func() (keys []string) {
		for _, key := range segment(cedict.Dawg, []rune(query)) {
			keys = append(keys, string(key))
		}
		return keys
	}())
	if len(results) == 1 && results[0]["m"] == "?" {
		log.Printf("q='%v' triggers Full-Text Search", query)
		results = keysToResults(searchDB(db, query))
	}
	if results == nil {
		log.Printf("q='%v' returns no results", query)
		results = []map[string]string{}
	}
	b, _ := json.Marshal(map[string]interface{}{"r": results})
	w.Write(b)
}

func keysToResults(keys []string) (results []map[string]string) {
	var m, p []string
	for _, key := range keys {
		entry, ok := cedict.Dict[key]
		if ok {
			m = entry.definitions
			p = entry.pinyin
		} else {
			m = []string{"?"}
			p = []string{""}
		}
		results = append(results, map[string]string{
			"z": key,
			"m": strings.Join(m, "/"),
			"p": strings.Join(p, "/"),
		})
	}
	return results
}

func serve(host, port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/segment", segmentHandler)
	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
