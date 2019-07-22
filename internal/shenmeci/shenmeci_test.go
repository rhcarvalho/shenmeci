package shenmeci

import (
	"testing"
)

func TestSegmentChinese(t *testing.T) {
	sentence := "语言信息处理"
	expectedWords := [][]rune{
		[]rune("语言"),
		[]rune("信息"),
		[]rune("处理"),
	}
	testSegment(sentence, expectedWords, t)
}

func TestNonChinese(t *testing.T) {
	sentence := "I am a sentence in English."
	expectedWords := [][]rune{
		[]rune("I am a sentence in English."),
	}
	testSegment(sentence, expectedWords, t)
}

func TestEnglishChinese(t *testing.T) {
	sentence := "Hello 世界."
	expectedWords := [][]rune{
		[]rune("Hello "),
		[]rune("世界"),
		[]rune("."),
	}
	testSegment(sentence, expectedWords, t)
}

func testSegment(sentence string, expectedWords [][]rune, t *testing.T) {
	d := cedict.Dawg
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
	for i := 0; i < b.N; i++ {
		LoadCEDICT()
	}
}

func BenchmarkSegment(b *testing.B) {
	d := cedict.Dawg
	sentence := []rune("语言信息处理English你好")
	for i := 0; i < b.N; i++ {
		segment(d, sentence)
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
