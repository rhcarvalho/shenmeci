package main

import (
	"testing"
)

func BenchmarkLoadCEDICT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadCEDICT("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz")
	}
}

func BenchmarkLoadCEDICT2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadCEDICT2("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz")
	}
}

func BenchmarkLoadCEDICT3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		loadCEDICT3("dict/cedict_1_0_ts_utf-8_mdbg.txt.gz")
	}
}
