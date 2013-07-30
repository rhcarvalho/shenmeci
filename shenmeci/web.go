package main

import (
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

type QueryRecord struct {
	Query       string
	Result      []map[string]string
	When        time.Time
	Duration    time.Duration
	RequestInfo *http.Request
}

func segmentHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	query := r.FormValue("q")
	results := keysToResults(func() (keys []string) {
		for _, key := range segment(cedict.Dawg, []rune(query)) {
			keys = append(keys, string(key))
		}
		return keys
	}())
	if len(results) == 1 && results[0]["m"] == "?" {
		log.Printf("q='%v' triggers Full-Text Search", query)
		results = keysToResults(searchDB(query))
	}
	if results == nil {
		log.Printf("q='%v' returns no results", query)
		results = []map[string]string{}
	}
	b, _ := json.Marshal(map[string]interface{}{"r": results})
	w.Write(b)
	duration := time.Since(startTime)

	// Insert into MongoDB in another goroutine.
	// This finishes the response without blocking.
	go collection.Insert(&QueryRecord{query, results, bson.Now(), duration, r})
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

var collection *mgo.Collection

func serve(host string, port int) {
	session, err := mgo.Dial(config.MongoURL)
	if err != nil {
		log.Fatal("MongoDB: ", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	collection = session.DB("shenmeci").C("queries")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(config.StaticPath, "index.html"))
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticPath))))
	http.HandleFunc("/segment", segmentHandler)
	log.Printf("serving at %s:%d", host, port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
