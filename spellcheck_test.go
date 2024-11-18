package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestCheckReader(t *testing.T) {
	wordList := []string{"abc", "def", "ghi", "jkl"}
	spellcheck := newSpellcheck()
	spellcheck.InitializeWordList(strings.NewReader(strings.Join(wordList, "\n")))
	textReader := strings.NewReader("Abc xxx DEF. \ngHi jkl zzz")
	spellingErrors := spellcheck.CheckReader(textReader)

	expectedLen := 2
	actualLen := 0
	for spellingError := range spellingErrors {
		t.Logf("\n%v\n", spellingError)
		actualLen++
	}
	if actualLen != expectedLen {

		t.Fatalf("\nExpected %v spelling errors; got %v\n", expectedLen, actualLen)
	}
}

func TestSuggestions(t *testing.T) {
	s := "abcdefg"
	trie := newTrieNode()
	wordList := []string{"a", "abx", "abcx", "abcdx", "abcdex", "abcdey", "abcdez"}
	trie.InsertAll(strings.NewReader(strings.Join(wordList, "\n")))

	expected := []string{"abcdex", "abcdey"}
	actual := getSuggestions(trie, s, 2)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected:\t%v\nActual:\t\t%v\n", expected, actual)
	}
}
func TestSuggestions2(t *testing.T) {
	s := "thes"
	trie := newTrieNode()
	wordList := []string{"these", "wraps", "are", "inky"}
	trie.InsertAll(strings.NewReader(strings.Join(wordList, "\n")))
	t.Logf("\ntrie:\n%s\n", trie)
	expected := []string{"these"}
	actual := getSuggestions(trie, s, 2)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected:\t%v\nActual:\t\t%v\n", expected, actual)
	}
}
