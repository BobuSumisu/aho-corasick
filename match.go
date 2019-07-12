package aho_corasick

import (
	"bytes"
	"fmt"
)

type Match struct {
	pos   int
	match []byte
}

func newMatch(pos int, match []byte) *Match {
	return &Match{pos: pos, match: match}
}

func newMatchString(pos int, match string) *Match {
	return &Match{pos: pos, match: []byte(match)}
}

func (m *Match) String() string {
	return fmt.Sprintf("{%d %q}", m.pos, m.match)
}

func (m *Match) Pos() int {
	return m.pos
}

func (m *Match) Match() []byte {
	return m.match
}

func (m *Match) MatchString() string {
	return string(m.match)
}

func MatchEqual(a, b *Match) bool {
	return a.pos == b.pos && bytes.Equal(a.match, b.match)
}
