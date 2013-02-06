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
	"path"
	"runtime"
	"unicode/utf8"
)

func segment(d *dawg.DAWG, sentence string) (words []string) {
	var longestWord string
	s := []rune(sentence)
	for len(s) > 0 {
		prefixes := d.Prefixes(string(s))
		if len(prefixes) > 0 {
			longestWord = prefixes[len(prefixes)-1]
		} else {
			longestWord = string(s[0])
		}
		words = append(words, longestWord)
		s = s[utf8.RuneCountInString(longestWord):]
	}
	return
}

func loadCEDICT(filename string) (map[string][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	br := bufio.NewReader(r)
	dict := make(map[string][]string)
	for {
		line, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if line[0] == '#' {
			continue
		}
		hanzi := string(bytes.Fields(line)[1])
		meaning := string(bytes.TrimSpace(line[bytes.IndexByte(line, '/'):]))
		dict[hanzi] = append(dict[hanzi], meaning)
		if err == io.EOF {
			break
		}
	}
	return dict, nil
}

// pathRelativeToExecutable returns an absolute path relative to
// this executable's directory.
func pathRelativeToExecutable(p string) string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), p)
}

func main() {
	filename := pathRelativeToExecutable("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz")
	dict, err := loadCEDICT(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	d := dawg.New(nil)
	for k := range dict {
		d.Insert(k)
	}
	fmt.Println(segment(d, "语言信息处理"))
}
