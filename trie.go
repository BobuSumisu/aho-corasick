package ahocorasick

const (
	rootState int = 1
	nilState  int = 0
)

// Trie represents a trie of patterns with extra links as per the Aho-Corasick algorithm.
type Trie struct {
	failLink []int
	dictLink []int
	dict     []int
	trans    [][256]int
}

type walkFn func(end, n int) bool

func (tr *Trie) walk(input []byte, fn walkFn) {
	s := rootState

	for i, c := range input {
		t := nilState

		if t = tr.trans[s][c]; t == nilState {
			for u := tr.failLink[s]; u != rootState; u = tr.failLink[u] {
				if t = tr.trans[u][c]; t != nilState {
					break
				}
			}

			if t == nilState {
				if t = tr.trans[rootState][c]; t == nilState {
					t = rootState
				}
			}
		}

		s = t

		if tr.dict[s] != 0 {
			if !fn(i, tr.dict[s]) {
				return
			}
		}

		if tr.dictLink[s] != nilState {
			for u := tr.dictLink[s]; u != nilState; u = tr.dictLink[u] {
				if !fn(i, tr.dict[u]) {
					return
				}
			}
		}
	}
}

// Match runs the Aho-Corasick string-search algorithm on a byte input.
func (tr *Trie) Match(input []byte) []*Match {
	matches := make([]*Match, 0)
	tr.walk(input, func(end, n int) bool {
		pos := end - n + 1
		matches = append(matches, &Match{pos: pos, match: input[pos : pos+n]})
		return true
	})
	return matches
}

// MatchFirst is the same as Match, but returns after first successful match.
func (tr *Trie) MatchFirst(input []byte) *Match {
	var match *Match
	tr.walk(input, func(end, n int) bool {
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
