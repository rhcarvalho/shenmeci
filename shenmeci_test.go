package main

import (
	"testing"
)

func BenchmarkLoadCEDICT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadCEDICT("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz")
	}
}
