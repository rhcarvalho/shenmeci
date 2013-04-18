package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func segmentHandler(w http.ResponseWriter, r *http.Request) {
	var m, p []string
	query := r.FormValue("q")
	// words is initialized to make sure its JSON representation is at
	// least an empty list and not null
	words := []interface{}{}
	for _, word := range segment(cedict.Dawg, []rune(query)) {
		z := string(word)
		entry, ok := cedict.Dict[z]
		if ok {
			m = entry.definitions
			p = entry.pinyin
		} else {
			m = []string{"?"}
			p = []string{""}
		}
		words = append(words, map[string]string{
			"z": z,
			"m": strings.Join(m, "/"),
			"p": strings.Join(p, "/"),
		})
	}
	b, _ := json.Marshal(map[string]interface{}{"r": words})
	w.Write(b)
}
