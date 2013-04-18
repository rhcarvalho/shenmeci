package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/rhcarvalho/DAWGo/dawg"
	"io"
	"log"
	"os"
	"os/exec"
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

func loadCEDICT(filename string) (c *CEDICT, err error) {
	if filename == "" {
		// if not set, try a default path
		filename = "dict/cedict_1_0_ts_utf-8_mdbg.txt.gz"
		if _, err := os.Stat(filename); err != nil {
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
	c = &CEDICT{Dict: make(map[string]CEDICTEntry), Dawg: dawg.New(nil)}
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
		pinyin := string(parts[2][bytes.IndexByte(parts[2], '[')+1 : bytes.IndexByte(parts[2], ']')])
		entry := c.Dict[wordSimplified]
		entry.definitions = append(entry.definitions, meaning)
		entry.pinyin = append(entry.pinyin, pinyinNumberedToUnicode(pinyin))
		c.Dict[wordSimplified] = entry
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
