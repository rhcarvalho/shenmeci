package dawg

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// the map keys are space-separated words to construct the DAWG.
var dawgs = map[string]*DAWG{
	"": {},
	"g": {children: map[rune]*DAWG{
		'g': {eow: true},
	}},
	"go": {children: map[rune]*DAWG{
		'g': {children: map[rune]*DAWG{
			'o': {eow: true},
		}},
	}},
	"g go": {children: map[rune]*DAWG{
		'g': {eow: true, children: map[rune]*DAWG{
			'o': {eow: true},
		}},
	}},
	"g t": {children: map[rune]*DAWG{
		'g': {eow: true},
		't': {eow: true},
	}},
	"go t": {children: map[rune]*DAWG{
		'g': {children: map[rune]*DAWG{
			'o': {eow: true},
		}},
		't': {eow: true},
	}},
	"语 语言 信 信息 处 处理": {children: map[rune]*DAWG{
		'处': {eow: true, children: map[rune]*DAWG{
			'理': {eow: true},
		}},
		'语': {eow: true, children: map[rune]*DAWG{
			'言': {eow: true},
		}},
		'信': {eow: true, children: map[rune]*DAWG{
			'息': {eow: true},
		}},
	}},
}

func TestNew(t *testing.T) {
	for words, d := range dawgs {
		nd := New(strings.Fields(words))
		if !dawgsEqual(nd, d) {
			t.Errorf("DAWG should be %v, got %v", d, nd)
		}
	}
}

func TestLongestCommonPrefix(t *testing.T) {
	type query struct {
		key           string
		longestPrefix string
	}
	tests := []struct {
		words   string
		queries []query
	}{
		{
			"g go", []query{
				{"", ""},
				{"g", "g"},
				{"go", "go"},
				{"golang", "go"},
				{"python", ""},
			},
		}, {
			"g t", []query{
				{"g", "g"},
				{"t", "t"},
				{"golang", "g"},
				{"tornado", "t"},
				{"z", ""},
			},
		}, {
			"", []query{
				{"", ""},
				{"g", ""},
				{"golang", ""},
			},
		}, {
			"语 语言 信 信息 处 处理", []query{
				{"语言信息处理", "语言"},
			},
		},
	}
	for _, test := range tests {
		d, ok := dawgs[test.words]
		if !ok {
			t.Errorf("Missing DAWG for words %#v", test.words)
			continue
		}
		for _, q := range test.queries {
			want := q.longestPrefix
			got := string(d.LongestCommonPrefix([]rune(q.key)))
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		}
	}
}

// Benchmarks

// Helper functions for printings DAWGs.
func (d *DAWG) String() string {
	eowS, chldS := eowToString(d.eow), childrenToString(d.children)
	var sep string
	if eowS != "" && chldS != "" {
		sep = ", "
	} else {
		sep = ""
	}
	return fmt.Sprintf("&DAWG{%s%s%s}", eowS, sep, chldS)
}

func eowToString(eow bool) string {
	if eow {
		return "eow: true"
	}
	return ""
}

func childrenToString(children map[rune]*DAWG) string {
	if children == nil {
		return ""
	}
	b := bytes.NewBufferString("children: map[rune]*DAWG{\n")
	for k, nd := range children {
		fmt.Fprintf(b, "'%q': %v,\n", k, nd)
	}
	b.WriteByte('}')
	return b.String()
}

// Helper functions for comparing DAWGs.
func dawgsEqual(d1, d2 *DAWG) bool {
	if d1 == nil {
		return d2 == nil
	}
	if d2 == nil {
		return false
	}
	if d1.eow != d2.eow {
		return false
	}
	for key, d1ChildNode := range d1.children {
		d2ChildNode, ok := d2.children[key]
		if !ok {
			return false
		}
		if !dawgsEqual(d1ChildNode, d2ChildNode) {
			return false
		}
	}
	return true
}
