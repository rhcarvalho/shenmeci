package pinyin

import "testing"

func TestPinyin(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"zhong1", "zhōng"},
		{"zhong1 guo2", "zhōng guó"},
		{"U S B shou3 zhi3", "U S B shǒu zhǐ"},
	}
	for _, test := range tests {
		if got := ToDiacritics(test.input); got != test.want {
			t.Errorf("got %q, want %q", got, test.want)
		}
	}
}
