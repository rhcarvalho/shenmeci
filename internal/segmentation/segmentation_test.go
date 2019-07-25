package segmentation

import (
	"strings"
	"testing"
)

var segmenter = NewSegmenter(strings.Fields("语言 信息 处理 世界"))

func TestSegment(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{"Empty", "", nil},
		{"Chinese", "语言信息处理", strings.Fields("语言 信息 处理")},
		{"NonChinese", "I am a sentence in English.", []string{"I am a sentence in English."}},
		{"EnglishChinese", "Hello 世界.", []string{"Hello ", "世界", "."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := segmenter.Segment([]rune(tt.text))
			if len(got) != len(tt.want) {
				t.Fatalf("got %d segments, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if string(got[i]) != tt.want[i] {
					t.Fatalf("[segment#%d] got %q, want %q", i+1, string(got[i]), tt.want[i])
				}
			}
		})
	}
}

func BenchmarkSegment(b *testing.B) {
	sentence := []rune("语言信息处理English你好")
	for i := 0; i < b.N; i++ {
		segmenter.Segment(sentence)
	}
}
