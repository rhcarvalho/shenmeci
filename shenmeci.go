package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/rhcarvalho/DAWGo/dawg"
	"io"
	"io/ioutil"
	"log"
	"os"
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

func fSegment(d *dawg.DAWG, r io.Reader, w io.Writer) {
	unsegmentedText, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range segment(d, bytes.Runes(unsegmentedText)) {
		// Do not print carriage returns
		if len(s) == 1 && s[0] == '\r' {
			continue
		}
		if i > 0 {
			w.Write([]byte{' '})
		}
		w.Write([]byte(string(s)))
	}
}

type CEDICT struct {
	Dict      map[string][]string
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
	c = &CEDICT{Dict: make(map[string][]string)}
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
	cedictPath      = flag.String("dict", os.Getenv("CEDICT"), "path to CEDICT")
	inFilePath      = flag.String("infile", "", "from where to read unsegmented text (defaults to stdin)")
	outFilePath     = flag.String("outfile", "", "where to write text segmented into words (defaults to stdout)")
	inFile, outFile *os.File
	err             error
)

func main() {
	flag.Parse()
	if *cedictPath == "" {
		log.Fatal("Missing environment variable CEDICT or command-line argument -dict.")
	}
	if *inFilePath == "" {
		inFile = os.Stdin
	} else {
		inFile, err = os.Open(*inFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()
	}
	if *outFilePath == "" {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create(*outFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
	}
	cedict, err := loadCEDICT(*cedictPath)
	if err != nil {
		log.Fatal(err)
	}
	d := dawg.New(nil)
	for k := range cedict.Dict {
		d.Insert(k)
	}
	fSegment(d, inFile, outFile)
}
