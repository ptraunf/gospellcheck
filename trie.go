package main

import (
	"bufio"
	"io"
	"slices"
	"strings"
)

type TrieNode struct {
	isKey    bool
	children map[rune]*TrieNode
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		isKey:    false,
		children: make(map[rune]*TrieNode),
	}
}

type Trie interface {
	Insert(key string) bool
	InsertAll(r io.Reader)
	Remove(key string) bool
	Contains(key string) bool
	LongestPrefix(key string) string
	KeysWithCommonPrefix(prefix string) []string
}

func (t *TrieNode) Contains(key string) bool {
	currentNode := t
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

func (t *TrieNode) addNewBranch(chars []rune) {
	currentNode := t

	for len(chars) >= 1 {
		c := chars[0]
		currentNode.children[c] = newTrieNode()
		currentNode = currentNode.children[c]
		chars = chars[1:]
	}
	currentNode.isKey = true

}
func (t *TrieNode) Insert(s string) bool {
	currentNode := t
	searchChars := []rune(s)
	for len(searchChars) >= 1 {
		c := searchChars[0]
		child, hasChild := currentNode.children[c]
		if hasChild {
			currentNode = child
			searchChars = searchChars[1:]
		} else {
			currentNode.addNewBranch(searchChars)
		}
	}
	currentNode.isKey = true
	return currentNode.isKey
}
func (t *TrieNode) Remove(key string) bool {
	return false
}
func (t *TrieNode) LongestPrefix(s string) string {
	if len(s) == 0 {
		return ""
	}
	currentNode := t
	chars := []rune(s)
	longestPrefix := make([]rune, 0)
	for len(chars) > 0 && len(currentNode.children) > 0 {
		c := chars[0]
		child, hasChild := currentNode.children[c]
		if hasChild {
			longestPrefix = append(longestPrefix, c)
			currentNode = child
			chars = chars[1:]
		} else {
			break
		}
	}
	return string(longestPrefix)
}

func (t *TrieNode) KeysWithCommonPrefix(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	currentNode := t
	var keys []string
	chars := []rune(s)
	longestPrefix := make([]rune, 0)
	for len(chars) > 0 && len(currentNode.children) > 0 {
		c := chars[0]
		child, hasChild := currentNode.children[c]
		if hasChild {
			longestPrefix = append(longestPrefix, c)
			currentNode = child
			chars = chars[1:]
		} else {
			keys = currentNode.Enumerate()
			break
		}
	}
	longestPrefixStr := string(longestPrefix)
	slices.Sort(keys)
	for i := 0; i < len(keys); i++ {
		keys[i] = longestPrefixStr + keys[i]
	}
	return keys
}

type NodePath struct {
	prefix string
	node   *TrieNode
}

func (t *TrieNode) Enumerate() []string {
	var neighborStack = make(stack[NodePath], 0)
	var enumeration []string
	var prefix = ""
	neighborStack.push(NodePath{prefix, t})

	for neighborStack.size() > 0 {
		currentPath := neighborStack.pop()
		if currentPath.node.isKey {
			enumeration = append(enumeration, currentPath.prefix)
		}
		for c, neighbor := range currentPath.node.children {
			neighborStack.push(
				NodePath{
					prefix: currentPath.prefix + string(c),
					node:   neighbor,
				})
		}
	}
	return enumeration
}

func (t *TrieNode) String() string {
	return strings.Join(t.Enumerate(), "\n")
}
func (t *TrieNode) InsertAll(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		word := scanner.Text()
		_ = t.Insert(word)
	}
}
