// Package dawg provides a Directed Acyclic Word Graph.
// A DAWG is a data structure optimized for fast string lookups.
package dawg

// A DAWG is the main type provided by this package.
type DAWG struct {
	children map[rune]*DAWG
	eow      bool // end-of-word marker
}

// New creates a new DAWG from a vocabulary.
func New(vocabulary []string) *DAWG {
	d := &DAWG{}
	for _, word := range vocabulary {
		d.Insert(word)
	}
	return d
}

// Insert a word into the DAWG.
func (d *DAWG) Insert(word string) {
	current := d
	for _, k := range word {
		if current.children == nil {
			current.children = make(map[rune]*DAWG)
		}
		if next, ok := current.children[k]; ok {
			current = next
		} else {
			next = &DAWG{}
			current.children[k] = next
			current = next
		}
	}
	current.eow = true
}

// LongestCommonPrefix returns the longest common prefix between sentence and
// all of the words inserted into the DAWG.
func (d *DAWG) LongestCommonPrefix(sentence []rune) []rune {
	current := d
	i, j := 0, 0
	for _, k := range sentence {
		if current.children == nil {
			break
		}
		if next, ok := current.children[k]; ok {
			i++
			if next.eow {
				j = i
			}
			current = next
		} else {
			break
		}
	}
	return sentence[:j]
}
