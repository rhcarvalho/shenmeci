package shenmeci

import (
	"testing"
)

func BenchmarkLoadCEDICT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LoadCEDICT()
	}
}

func TestPinyin(t *testing.T) {
	if r := pinyinNumberedSyllableToUnicode("zhong1"); r != "zhōng" {
		t.Errorf("expected %v got %v", "zhōng", r)
	}
	if r := pinyinNumberedToUnicode("zhong1 guo2"); r != "zhōng guó" {
		t.Errorf("expected %v got %v", "zhōng guó", r)
	}
	if r := pinyinNumberedToUnicode("U S B shou3 zhi3"); r != "U S B shǒu zhǐ" {
		t.Errorf("expected %v got %v", "U S B shǒu zhǐ", r)
	}
}

func init() {
	LoadConfig()
	ValidateConfig()
	LoadCEDICT()
}
