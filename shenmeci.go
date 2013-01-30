package main

import (
	"fmt"
	"github.com/rhcarvalho/DAWGo/dawg"
)

func segment(d *dawg.DAWG, sentence string) (words []string) {
	var longestWord string
	for len(sentence) > 0 {
		prefixes := d.Prefixes(sentence)
		if len(prefixes) > 0 {
			longestWord = prefixes[len(prefixes)-1]
		} else {
			longestWord = sentence[:1]
		}
		words = append(words, longestWord)
		sentence = sentence[len(longestWord):]
	}
	return
}

func main() {
	d := dawg.New([]string{"go", "python", "ruby", "c", "cpp"})
	fmt.Println(segment(d, "golangpythoncpp"))
}
