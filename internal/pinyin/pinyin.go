package pinyin

import (
	"strconv"
	"strings"
)

// ToDiacritics converts a string of Pinyin syllables to use diacritics instead
// of numbers.
func ToDiacritics(pinyin string) string {
	syllables := strings.Fields(pinyin)
	for i, syllable := range syllables {
		syllables[i] = toDiacritics(syllable)
	}
	return strings.Join(syllables, " ")
}

// toDiacritics converts a single Pinyin syllable from numbered form to use
// diacritics.
// Example: zhong1 -> zhōng
func toDiacritics(pinyin string) string {
	var (
		a = []rune("āáǎàa")
		e = []rune("ēéěèe")
		i = []rune("īíǐìi")
		o = []rune("ōóǒòo")
		u = []rune("ūúǔùu")
		v = []rune("ǖǘǚǜü")
	)
	if len := len(pinyin); len > 0 {
		last := pinyin[len-1:]
		tone, err := strconv.Atoi(last)
		if err == nil && 1 <= tone && tone <= 5 {
			pinyin = pinyin[:len-1]
			switch {
			case strings.Contains(pinyin, "a"):
				pinyin = strings.Replace(pinyin, "a", string(a[tone-1]), 1)
			case strings.Contains(pinyin, "e"):
				pinyin = strings.Replace(pinyin, "e", string(e[tone-1]), 1)
			case strings.Contains(pinyin, "ou"):
				pinyin = strings.Replace(pinyin, "o", string(o[tone-1]), 1)
			case strings.Contains(pinyin, "io"):
				pinyin = strings.Replace(pinyin, "o", string(o[tone-1]), 1)
			case strings.Contains(pinyin, "iu"):
				pinyin = strings.Replace(pinyin, "u", string(u[tone-1]), 1)
			case strings.Contains(pinyin, "ui"):
				pinyin = strings.Replace(pinyin, "i", string(i[tone-1]), 1)
			case strings.Contains(pinyin, "uo"):
				pinyin = strings.Replace(pinyin, "o", string(o[tone-1]), 1)
			case strings.Contains(pinyin, "i"):
				pinyin = strings.Replace(pinyin, "i", string(i[tone-1]), 1)
			case strings.Contains(pinyin, "o"):
				pinyin = strings.Replace(pinyin, "o", string(o[tone-1]), 1)
			case strings.Contains(pinyin, "u:"):
				pinyin = strings.Replace(pinyin, "u:", string(v[tone-1]), 1)
			case strings.Contains(pinyin, "u"):
				pinyin = strings.Replace(pinyin, "u", string(u[tone-1]), 1)
			}
		}
		// Make sure there is no "u:" left.
		// Example: "lu:e4" => "lüè".
		if strings.Contains(pinyin, "u:") {
			pinyin = strings.Replace(pinyin, "u:", "ü", 1)
		}
	}
	return pinyin
}
