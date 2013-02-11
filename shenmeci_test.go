package main

import (
	"github.com/rhcarvalho/DAWGo/dawg"
	"os"
	"testing"
)

func TestSegment(t *testing.T) {
	d, err := newDAWGFromCEDICT()
	if err != nil {
		t.Fatal(err)
	}
	sentence := "语言信息处理"
	expectedWords := [][]rune{
		[]rune("语言"),
		[]rune("信息"),
		[]rune("处理"),
	}
	words := segment(d, []rune(sentence))
	if len(words) != len(expectedWords) {
		t.Errorf("segmented %q should be %q, got %q", sentence, expectedWords, words)
	}
	for i, word := range words {
		if string(word) != string(expectedWords[i]) {
			t.Errorf("segmented %q should be %q, got %q\nfirst differing word [%d]: %q != %q",
				sentence, expectedWords, words, i, expectedWords[i], word)
		}
	}
}

func BenchmarkLoadCEDICT(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if _, err = loadCEDICT(os.Getenv("CEDICT")); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSegment(b *testing.B) {
	b.StopTimer()
	d, err := newDAWGFromCEDICT()
	if err != nil {
		b.Fatal(err)
	}
	sentence := []rune("语言信息处理")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		segment(d, sentence)
	}
}

func newDAWGFromCEDICT() (d *dawg.DAWG, err error) {
	cedict, err := loadCEDICT(os.Getenv("CEDICT"))
	if err != nil {
		return nil, err
	}
	d = dawg.New(nil)
	for k := range cedict.Dict {
		d.Insert(k)
	}
	return
}