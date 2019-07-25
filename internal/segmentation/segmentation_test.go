package segmentation

import (
	"strings"
	"testing"
)

var segmenter = NewSegmenter(strings.Fields("语言 信息 处理 世界"))

func TestSegmentEmpty(t *testing.T) {
	testSegment("", [][]rune{}, t)
}

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
	words := segmenter.Segment([]rune(sentence))
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

func BenchmarkSegment(b *testing.B) {
	sentence := []rune("语言信息处理English你好")
	for i := 0; i < b.N; i++ {
		segmenter.Segment(sentence)
	}
}
