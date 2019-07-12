package aho_corasick

type trieNode struct {
	parent      *trieNode
	fail        *trieNode
	dict        *trieNode
	c           byte
	patternLen  int
	transitions [256]*trieNode
}

func newTrieNode(parent *trieNode, c byte) *trieNode {
	return &trieNode{
		parent:     parent,
		fail:       nil,
		dict:       nil,
		c:          c,
		patternLen: 0,
	}
}

type Trie struct {
	root *trieNode
}

type walkFn func(end, n int) bool

func (tr *Trie) walk(input []byte, fn walkFn) {
	node := tr.root
	var next *trieNode

	for i, c := range input {

		if next = node.transitions[c]; next == nil {
			for n := node.fail; n != tr.root; n = n.fail {
				if next = n.transitions[c]; next != nil {
					break
				}
			}

			if next == nil {
				if next = tr.root.transitions[c]; next == nil {
					next = tr.root
				}
			}
		}

		node = next

		if node.patternLen != 0 {
			if !fn(i, node.patternLen) {
				return
			}
		}

		if node.dict != nil {
			for n := node.dict; n != nil; n = n.dict {
				if !fn(i, n.patternLen) {
					return
				}
			}
		}
	}
}

func (tr *Trie) Match(input []byte) []*Match {
	matches := make([]*Match, 0)
	tr.walk(input, func(end, n int) bool {
		pos := end - n + 1
		matches = append(matches, &Match{pos: pos, match: input[pos : pos+n]})
		return true
	})
	return matches
}

func (tr *Trie) MatchFirst(input []byte) *Match {
	var match *Match
	tr.walk(input, func(end, n int) bool {
		pos := end - n + 1
		match = &Match{pos: pos, match: input[pos : pos+n]}
		return false
	})
	return match
}

func (tr *Trie) MatchString(input string) []*Match {
	return tr.Match([]byte(input))
}

func (tr *Trie) MatchFirstString(input string) *Match {
	return tr.MatchFirst([]byte(input))
}
