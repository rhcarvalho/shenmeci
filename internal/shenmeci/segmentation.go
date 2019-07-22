package shenmeci

import (
	"github.com/rhcarvalho/shenmeci/internal/segmentation/dawg"
)

func segment(d *dawg.DAWG, sentence []rune) (words [][]rune) {
	var nextWord []rune
	for len(sentence) > 0 {
		nextWord = longestPrefixWord(d, sentence)
		words = append(words, nextWord)
		sentence = sentence[len(nextWord):]
	}
	return
}

func longestPrefixWord(d *dawg.DAWG, sentence []rune) (word []rune) {
	prefixes := d.Prefixes(sentence)
	if len(prefixes) > 0 {
		word = prefixes[len(prefixes)-1]
	} else {
		// In this case, there is no word in the dictionary that is a
		// prefix of the sentence, so we take the longest non-prefix
		// portion of the sentence as the longest prefix.
		// This means that unknown terms are not segmented.
		for len(sentence) > 0 {
			prefixes := d.Prefixes(sentence)
			if len(prefixes) > 0 {
				break
			}
			word = append(word, sentence[0])
			sentence = sentence[1:]
		}
	}
	return
}

func segmentNew(d *dawg.DAWG, sentence []rune) (words [][]rune) {
	var nextWord []rune
	for len(sentence) > 0 {
		nextWord = longestPrefixWordNew(d, sentence)
		words = append(words, nextWord)
		sentence = sentence[len(nextWord):]
	}
	return
}

func longestPrefixWordNew(d *dawg.DAWG, sentence []rune) (word []rune) {
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
