package ahocorasick

import (
	"bytes"
	"fmt"
)

// Match represents a matched pattern in the input.
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

// Pos returns the byte position of the match.
func (m *Match) Pos() int {
	return m.pos
}

// Match returns the pattern matched.
func (m *Match) Match() []byte {
	return m.match
}

// MatchString returns the pattern matched as a string.
func (m *Match) MatchString() string {
	return string(m.match)
}

// MatchEqual check whether two matches are equal (i.e. at same position and same pattern).
func MatchEqual(a, b *Match) bool {
	return a.pos == b.pos && bytes.Equal(a.match, b.match)
}
