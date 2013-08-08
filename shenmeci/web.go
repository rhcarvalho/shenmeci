package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type QueryRecord struct {
	Query       string
	Result      []*Result
	When        time.Time
	Duration    time.Duration
	RequestInfo *http.Request
}

type Results struct {
	R []*Result
}

type Result struct {
	Z string // Hanzi
	M string // Meaning
	P string // Pinyin
}

func (r *Results) MarshalJSON() ([]byte, error) {
	if r.R != nil {
		return json.Marshal(&map[string]interface{}{"r": r.R})
	} else {
		return json.Marshal(&map[string]interface{}{"r": []interface{}{}})
	}
}

func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&map[string]string{
			"z": template.HTMLEscapeString(r.Z),
			"m": template.HTMLEscapeString(r.M),
			"p": template.HTMLEscapeString(r.P)})
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
	if len(results.R) == 1 && results.R[0].M == "?" {
		log.Printf("q='%v' triggers Full-Text Search", query)
		results = keysToResults(searchDB(query))
	}
	if len(results.R) == 0 {
		log.Printf("q='%v' returns no results", query)
	}
	b, _ := json.Marshal(results)
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	duration := time.Since(startTime)

	// Insert into MongoDB in another goroutine.
	// This finishes the response without blocking.
	go func() {
		err := collection.Insert(&QueryRecord{query, results.R, startTime, duration, r})
		// Log and refresh the Session in case of insertion errors
		if err != nil {
			log.Print("MongoDB: ", err)
			collection.Database.Session.Refresh()
		}
	}()
}

func keysToResults(keys []string) *Results {
	results := &Results{}
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
		results.R = append(results.R, &Result{
			Z: key,
			M: strings.Join(m, "/"),
			P: strings.Join(p, "/"),
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
