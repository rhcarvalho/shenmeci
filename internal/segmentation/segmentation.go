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
		segment = s.nextSegment(text)
		segments = append(segments, segment)
		text = text[len(segment):]
	}
	return segments
}

// nextSegment returns the first segment in text. A segment is always a prefix
// of text, and it is either the longest word in the Segmenter vocabulary that
// shares a prefix with text or the longest sequence of non-words.
func (s *dawgSegmenter) nextSegment(text []rune) (segment []rune) {
	d := s.dawg
	segment = d.LongestCommonPrefix(text)
	if len(segment) > 0 {
		return segment
	}
	// There is no word in the vocabulary that is a prefix of the text,
	// return the longest non-word portion of the text.
	// This means that consecutive unknown terms become a single segment.
	// In practice, for a Segmenter created with a vocabulary of Chinese
	// terms, this will keep numbers and English text as a single segment.
	for len(text) > 0 && len(d.LongestCommonPrefix(text)) == 0 {
		segment = append(segment, text[0])
		text = text[1:]
	}
	return segment
}
