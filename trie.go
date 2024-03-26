package gospellcheck

import "fmt"

type Node struct {
	isKey    bool
	children map[byte]Node
}

type Trie struct {
	root Node
}

func (t Trie) New() *Trie {
	trie := Trie{
		root: Node{
			isKey:    false,
			children: make(map[byte]Node),
		},
	}
	return &trie
}

func (t Trie) Contains(key string) bool {
	if "" == key {
		return t.root.isKey
	}
	for i, b := range key {
		fmt.Printf("i: %v, b: %v\n", i, b)
	}
	return false
}
func (t *Trie) Insert(s string) bool {
	return false
}
func (t *Trie) Remove(key string) bool {
	return false
}
func (t *Trie) LongestPrefix(key string) string {
	return ""
}
func (t *Trie) KeysStartingWith(prefix string) []string {
	return []string{}
}

type StringContainer interface {
	Insert(key string) bool
	Remove(key string) bool
	Contains(key string) bool
	LongestPrefix(key string) string
	KeysStartingWith(prefix string) []string
}
