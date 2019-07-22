package segmentation

import (
	"github.com/rhcarvalho/shenmeci/internal/segmentation/dawg"
)

// A Segmenter segments text.
type Segmenter interface {
	Segment(text []rune) (segments [][]rune)
}

type dawgSegmenter struct {
	dawg *dawg.DAWG
}

// NewSegmenter returns a new Segmenter that knows about words in vocabulary.
// Each segment is either a word from vocabulary or the longest sequence of
// consecutive non-words.
func NewSegmenter(vocabulary []string) Segmenter {
	return &dawgSegmenter{
		dawg: dawg.New(vocabulary),
	}
}

// Segment splits text into segments of words and non-words according to its
// vocabulary.
func (s *dawgSegmenter) Segment(text []rune) (segments [][]rune) {
	var segment []rune
	for len(text) > 0 {
		segment = s.longestPrefixWord(text)
		segments = append(segments, segment)
		text = text[len(segment):]
	}
	return segments
}

func (s *dawgSegmenter) longestPrefixWord(sentence []rune) (word []rune) {
	d := s.dawg
	longestPrefix := d.LongestCommonPrefix(sentence)
	if len(longestPrefix) > 0 {
		return longestPrefix
	}
	// In this case, there is no word in the dictionary that is a
	// prefix of the sentence, so we take the longest non-prefix
	// portion of the sentence as the longest prefix.
	// This means that unknown terms are not segmented.
	for len(sentence) > 0 {
		longestPrefix = d.LongestCommonPrefix(sentence)
		if len(longestPrefix) > 0 {
			break
		}
		word = append(word, sentence[0])
		sentence = sentence[1:]
	}
	return
}
