package main

import (
	"github.com/rhcarvalho/DAWGo/dawg"
	"testing"
)

func BenchmarkLoadCEDICT(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if _, err = loadCEDICT("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSegment(b *testing.B) {
	b.StopTimer()
	cedict, err := loadCEDICT("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz")
	if err != nil {
		b.Fatal(err)
	}
	d := dawg.New(nil)
	for k := range cedict.Dict {
		d.Insert(k)
	}
	sentence := []rune("语言信息处理")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		segment(d, sentence)
	}
}
