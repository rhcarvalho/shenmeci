package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/rhcarvalho/DAWGo/dawg"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func segment(d *dawg.DAWG, sentence []rune) (words [][]rune) {
	var nextWord []rune
	for len(sentence) > 0 {
		nextWord = longestPrefixWord(d, sentence)
		words = append(words, nextWord)
		sentence = sentence[len(nextWord):]
	}
	return
}

func longestPrefixWord(d *dawg.DAWG, sentence []rune) (word []rune) {
	prefixes := d.Prefixes(sentence)
	if len(prefixes) > 0 {
		word = prefixes[len(prefixes)-1]
	} else {
		word = sentence[:1]
	}
	return
}

func segmentHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	// words is initialized to make sure its JSON representation is at
	// least an empty list and not null
	words := []interface{}{}
	for _, word := range segment(cedict.Dawg, []rune(query)) {
		z := string(word)
		m, ok := cedict.Dict[z]
		if !ok {
			m = []string{"?"}
		}
		words = append(words, map[string]interface{}{"z": z, "m": strings.Join(m, "/")})
	}
	b, _ := json.Marshal(map[string]interface{}{"r": words})
	w.Write(b)
}

type CEDICT struct {
	Dict      map[string][]string
	Dawg      *dawg.DAWG
	MaxKeyLen int
}

func loadCEDICT(filename string) (c *CEDICT, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return
	}
	defer r.Close()
	br := bufio.NewReader(r)
	c = &CEDICT{Dict: make(map[string][]string), Dawg: dawg.New(nil)}
	for {
		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if line[0] == '#' {
			continue
		}
		// The basic format of a CC-CEDICT entry is:
		//   Traditional Simplified [pin1 yin1] /English equivalent 1/equivalent 2/
		// For example:
		//   中國 中国 [Zhong1 guo2] /China/Middle Kingdom/
		parts := bytes.SplitN(line, []byte{' '}, 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("line not in CEDICT format: %q", line)
		}
		wordSimplified := string(parts[1])
		meaning := string(bytes.TrimSpace(parts[2][bytes.IndexByte(parts[2], '/'):]))
		c.Dict[wordSimplified] = append(c.Dict[wordSimplified], meaning)
		c.Dawg.Insert(wordSimplified)
		if wordLen := len([]rune(wordSimplified)); wordLen > c.MaxKeyLen {
			c.MaxKeyLen = wordLen
		}
		if err == io.EOF {
			break
		}
	}
	return
}

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
				// file does not exist
				log.Fatal("Missing environment variable CEDICT.")
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
