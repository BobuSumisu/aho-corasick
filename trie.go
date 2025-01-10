package ahocorasick

import (
	"sync"
)

const (
	rootState uint32 = 1
	nilState  uint32 = 0
)

// Trie represents a trie of patterns with extra links as per the Aho-Corasick algorithm.
type Trie struct {
	failTrans [][256]uint32

	dict     []uint32
	pattern  []uint32
	dictLink []uint32

	matchPool       sync.Pool // Pool for match slice pointers
	matchStructPool sync.Pool // Pool for Match structs
}

// Walk calls this function on any match, giving the end position, length of the matched bytes,
// and the pattern number.
type WalkFn func(end, n, pattern uint32) bool

// Walk runs the algorithm on a given output, calling the supplied callback function on every
// match. The algorithm will terminate if the callback function returns false.
func (tr *Trie) Walk(input []byte, fn WalkFn) {
	// Local references to frequently accessed slices.
	failTrans := tr.failTrans
	dict := tr.dict
	pattern := tr.pattern
	dictLink := tr.dictLink

	s := rootState

	inputLen := len(input)
	for i := range inputLen {
		s = failTrans[s][input[i]]

		ds := dict[s]
		dl := dictLink[s]
		if ds != 0 || dl != nilState {
			if ds != 0 && !fn(uint32(i), ds, pattern[s]) {
				return
			}
			for u := dl; u != nilState; u = dictLink[u] {
				if !fn(uint32(i), dict[u], pattern[u]) {
					return
				}
			}
		}
	}
}

// Match runs the Aho-Corasick string-search algorithm on a byte input.
func (tr *Trie) Match(input []byte) []*Match {
	matches := tr.matchPool.Get().(*[]*Match)
	*matches = (*matches)[:0] // Reset slice while keeping capacity

	tr.Walk(input, func(end, n, pattern uint32) bool {
		pos := end - n + 1
		match := tr.matchStructPool.Get().(*Match)
		match.pos = pos
		match.pattern = pattern
		match.match = input[pos : pos+n]
		*matches = append(*matches, match)
		return true
	})

	result := make([]*Match, len(*matches))
	copy(result, *matches)

	// Return matches to pool
	for i := range *matches {
		(*matches)[i] = nil
	}
	tr.matchPool.Put(matches)

	return result
}

// MatchFirst is the same as Match, but returns after first successful match.
func (tr *Trie) MatchFirst(input []byte) *Match {
	var match *Match
	tr.Walk(input, func(end, n, pattern uint32) bool {
		pos := end - n + 1
		match = &Match{pos: pos, match: input[pos : pos+n]}
		return false
	})
	return match
}

// MatchString runs the Aho-Corasick string-search algorithm on a string input.
func (tr *Trie) MatchString(input string) []*Match {
	return tr.Match([]byte(input))
}

// MatchFirstString is the same as MatchString, but returns after first successful match.
func (tr *Trie) MatchFirstString(input string) *Match {
	return tr.MatchFirst([]byte(input))
}
