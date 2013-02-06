package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"github.com/rhcarvalho/DAWGo/dawg"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func segment(d *dawg.DAWG, sentence []rune) (words [][]rune) {
	var longestWord []rune
	for len(sentence) > 0 {
		prefixes := d.Prefixes(sentence)
		if len(prefixes) > 0 {
			longestWord = prefixes[len(prefixes)-1]
		} else {
			longestWord = sentence[:1]
		}
		words = append(words, longestWord)
		sentence = sentence[len(longestWord):]
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
	dict, err := loadCEDICT(*cedictPath)
	if err != nil {
		log.Fatal(err)
	}
	d := dawg.New(nil)
	for k := range dict {
		d.Insert(k)
	}
	unsegmentedText, err := ioutil.ReadAll(inFile)
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range segment(d, bytes.Runes(unsegmentedText)) {
		// Do not print carriage returns
		if len(s) == 1 && s[0] == '\r' {
			continue
		}
		if i > 0 {
			outFile.Write([]byte{' '})
		}
		outFile.WriteString(string(s))
	}
}
