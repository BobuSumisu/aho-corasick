package aho_corasick

type TrieBuilder struct {
	root *trieNode
}

func NewTrieBuilder() *TrieBuilder {
	return &TrieBuilder{
		root: newTrieNode(nil, 0),
	}
}

func (tb *TrieBuilder) AddPattern(pattern []byte) *TrieBuilder {
	node := tb.root
	var next *trieNode

	for _, c := range pattern {

		if next = node.transitions[c]; next == nil {
			next = newTrieNode(node, c)
			node.transitions[c] = next
		}
		node = next
	}

	node.patternLen = len(pattern)

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
	return &Trie{root: tb.root}
}

func (tb *TrieBuilder) computeFailLinks(node *trieNode) {
	if node.fail != nil {
		return
	}

	if node == tb.root || node.parent == tb.root {
		node.fail = tb.root
	} else {
		for n := node.parent.fail; n != tb.root; n = n.fail {
			if n.fail == nil {
				tb.computeFailLinks(n)
			}

			if node.fail = n.transitions[node.c]; node.fail != nil {
				break
			}
		}

		if node.fail == nil {
			if node.fail = tb.root.transitions[node.c]; node.fail == nil {
				node.fail = tb.root
			}
		}
	}

	for _, child := range node.transitions {
		if child != nil {
			tb.computeFailLinks(child)
		}
	}
}

func (tb *TrieBuilder) computeDictLinks(node *trieNode) {
	if node != tb.root {
		for n := node.fail; n != tb.root; n = n.fail {
			if n.patternLen != 0 {
				node.dict = n
				break
			}
		}
	}

	for _, child := range node.transitions {
		if child != nil {
			tb.computeDictLinks(child)
		}
	}
}
