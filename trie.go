package main

import (
	// "fmt"
	"bufio"
	"log"
	"os"
	"strings"
)

type Node struct {
	isKey    bool
	children map[rune]*Node
}

func newNode() *Node {
	return &Node{
		isKey:    false,
		children: make(map[rune]*Node),
	}
}

type Trie struct {
	root *Node
}

func newTrie() *Trie {
	trie := Trie{
		root: newNode(),
	}
	return &trie
}

func (t Trie) Contains(key string) bool {
	currentNode := t.root
	searchChars := []rune(key)
	for len(searchChars) >= 1 {
		c := searchChars[0]
		child, hasChild := currentNode.children[c]
		if hasChild {
			currentNode = child
			searchChars = searchChars[1:]
		} else {
			return false
		}
	}
	return currentNode.isKey
}

func (t *Trie) addNewBranch(node *Node, chars []rune) {
	currentNode := node

	for len(chars) >= 1 {
		c := chars[0]
		currentNode.children[c] = newNode()
		currentNode = currentNode.children[c]
		chars = chars[1:]
	}
	currentNode.isKey = true

}
func (t *Trie) Insert(s string) bool {
	currentNode := t.root
	searchChars := []rune(s)
	for len(searchChars) >= 1 {
		c := searchChars[0]
		child, hasChild := currentNode.children[c]
		if hasChild {
			currentNode = child
			searchChars = searchChars[1:]
		} else {
			t.addNewBranch(currentNode, searchChars)
		}
	}
	currentNode.isKey = true
	return currentNode.isKey
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

type stack[T any] struct {
	data []T
}

func (s *stack[T]) push(item T) {
	s.data = append(s.data, item)
}
func (s *stack[T]) pop() T {
	n := len(s.data) - 1
	item := s.data[n]
	s.data = s.data[:n]
	return item
}
func (s *stack[T]) peek() T {
	n := len(s.data) - 1
	return s.data[n]
}

func (s *stack[T]) size() int {
	return len(s.data)
}

type NodePath struct {
	prefix string
	node   *Node
}

func enumerate(n *Node) []string {
	var neighborStack stack[NodePath]
	var enumeration []string
	var prefix = ""
	neighborStack.push(NodePath{prefix, n})

	for neighborStack.size() > 0 {
		currentPath := neighborStack.pop()
		if currentPath.node.isKey {
			// fmt.Println(currentPath.prefix)
			enumeration = append(enumeration, currentPath.prefix)
		}
		for c, neigbor := range currentPath.node.children {
			neighborStack.push(
				NodePath{
					prefix: currentPath.prefix + string(c),
					node:   neigbor,
				})
		}
	}
	return enumeration
}
func (t *Trie) Enumerate() []string {
	return enumerate(t.root)
}

func (t *Trie) String() string {
	return strings.Join(enumerate(t.root), "\n")
}

type StringContainer interface {
	Insert(key string) bool
	Remove(key string) bool
	Contains(key string) bool
	LongestPrefix(key string) string
	KeysStartingWith(prefix string) []string
}

func initializeSpellcheck(filename string) *Trie {
	var trie = newTrie()
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := scanner.Text()
		trie.Insert(word)
	}

	return trie
}
