package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"os"

	"github.com/rhcarvalho/DAWGo/dawg"
)

type CEDICTEntry struct {
	definitions []string
	pinyin      []string
}

type CEDICT struct {
	Dict      map[string]CEDICTEntry
	Dawg      *dawg.DAWG
	MaxKeyLen int
}

var cedict *CEDICT

func loadCEDICT() {
	f, err := os.Open(config.CedictPath)
	if err != nil {
		log.Fatal("CEDICT: ", err)
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal("CEDICT: ", err)
	}
	defer r.Close()
	br := bufio.NewReader(r)
	cedict = &CEDICT{Dict: make(map[string]CEDICTEntry), Dawg: dawg.New(nil)}
	log.Println("loading CEDICT into dict/DAWG...")
	for {
		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal("CEDICT: ", err)
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
			log.Fatalf("line not in CEDICT format: %q", line)
		}
		wordSimplified := string(parts[1])
		meaning := string(bytes.TrimSpace(parts[2][bytes.IndexByte(parts[2], '/'):]))
		pinyin := string(parts[2][bytes.IndexByte(parts[2], '[')+1 : bytes.IndexByte(parts[2], ']')])
		entry := cedict.Dict[wordSimplified]
		entry.definitions = append(entry.definitions, meaning)
		entry.pinyin = append(entry.pinyin, pinyinNumberedToUnicode(pinyin))
		cedict.Dict[wordSimplified] = entry
		cedict.Dawg.Insert(wordSimplified)
		if wordLen := len([]rune(wordSimplified)); wordLen > cedict.MaxKeyLen {
			cedict.MaxKeyLen = wordLen
		}
		if err == io.EOF {
			break
		}
	}
	log.Printf("loaded %v entries\n", len(cedict.Dict))
}
