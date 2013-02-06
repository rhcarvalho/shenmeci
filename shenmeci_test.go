package main

import (
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
