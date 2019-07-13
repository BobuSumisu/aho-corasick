package aho_corasick

type state struct {
	id       int
	value    byte
	parent   *state
	failLink *state
	dictLink *state
	dict     int
	trans    map[byte]*state
}

func newState(id int, value byte, parent *state) *state {
	return &state{
		id:     id,
		value:  value,
		parent: parent,
		trans:  make(map[byte]*state),
	}
}

type TrieBuilder struct {
	states []*state
	root   *state
}

func NewTrieBuilder() *TrieBuilder {
	tb := &TrieBuilder{
		states: make([]*state, 0),
	}
	tb.addState(0, nil)
	tb.addState(0, nil)
	tb.root = tb.states[1]
	return tb
}

func (tb *TrieBuilder) addState(value byte, parent *state) *state {
	s := newState(len(tb.states), value, parent)
	tb.states = append(tb.states, s)
	return s
}

func (tb *TrieBuilder) AddPattern(pattern []byte) *TrieBuilder {
	s := tb.root
	var t *state
	var found bool

	for _, c := range pattern {
		if t, found = s.trans[c]; !found {
			t = tb.addState(c, s)
			s.trans[c] = t
		}
		s = t
	}

	s.dict = len(pattern)

	return tb
}

func (tb *TrieBuilder) AddPatterns(patterns [][]byte) *TrieBuilder {
	for _, pattern := range patterns {
		tb.AddPattern(pattern)
	}
	return tb
}

func (tb *TrieBuilder) AddString(pattern string) *TrieBuilder {
	return tb.AddPattern([]byte(pattern))
}

func (tb *TrieBuilder) AddStrings(patterns []string) *TrieBuilder {
	for _, pattern := range patterns {
		tb.AddString(pattern)
	}
	return tb
}

func (tb *TrieBuilder) Build() *Trie {
	tb.computeFailLinks(tb.root)
	tb.computeDictLinks(tb.root)

	numStates := len(tb.states)

	parent := make([]int, numStates)
	failLink := make([]int, numStates)
	dictLink := make([]int, numStates)
	dict := make([]int, numStates)
	trans := make([][256]int, numStates)

	for i, s := range tb.states {
		if i == 0 {
			continue
		}
		if s.parent != nil {
			parent[i] = s.parent.id
		}
		failLink[i] = s.failLink.id
		if s.dictLink != nil {
			dictLink[i] = s.dictLink.id
		}
		dict[i] = s.dict
		for c, t := range s.trans {
			trans[i][c] = t.id
		}
	}

	return &Trie{
		failLink: failLink,
		dictLink: dictLink,
		dict:     dict,
		trans:    trans,
	}
}

func (tb *TrieBuilder) computeFailLinks(s *state) {
	if s.failLink != nil {
		return
	}

	if s == tb.root || s.parent == tb.root {
		s.failLink = tb.root
	} else {
		for t := s.parent.failLink; t != tb.root; t = t.failLink {
			if t.failLink == nil {
				tb.computeFailLinks(t)
			}

			if s.failLink = t.trans[s.value]; s.failLink != nil {
				break
			}
		}

		if s.failLink == nil {
			if s.failLink = tb.root.trans[s.value]; s.failLink == nil {
				s.failLink = tb.root
			}
		}
	}

	for _, child := range s.trans {
		tb.computeFailLinks(child)
	}
}

func (tb *TrieBuilder) computeDictLinks(s *state) {
	if s != tb.root {
		for t := s.failLink; t != tb.root; t = t.failLink {
			if t.dict != 0 {
				s.dictLink = t
				break
			}
		}
	}

	for _, child := range s.trans {
		tb.computeDictLinks(child)
	}
}
