package main

import (
	"strings"
	"testing"
)

func TestCheckLine(t *testing.T) {
	dictionary := newTrieNode()
	wordList := []string{"abc", "def", "ghi", "jkl"}
	for _, word := range wordList {
		dictionary.Insert(word)
	}
	text := "Abc xxx def yyy. Ghi zzz"
	spellingErrors := checkLine(dictionary, text, 1)

	if len(spellingErrors) != 3 {
		t.Fatalf("\nExpected 3 spelling errors; got %v\n", len(spellingErrors))
	}
	for _, spellingError := range spellingErrors {
		t.Logf("\n%s\n", spellingError)
	}
}

func TestCheckReaderSequential(t *testing.T) {
	dictionary := newTrieNode()
	wordList := []string{"abc", "def", "ghi", "jkl"}
	for _, word := range wordList {
		dictionary.Insert(word)
	}
	textReader := strings.NewReader("Abc xxx DEF. \ngHi jkl zzz")

	spellingErrors := checkReaderSequential(dictionary, textReader)

	if len(spellingErrors) != 2 {

		t.Fatalf("\nExpected 2 spelling errors; got %v\n", len(spellingErrors))
	}
	for _, spellingError := range spellingErrors {
		t.Logf("\n%v\n", spellingError)
	}
}
func TestCheckReaderConcurrent(t *testing.T) {
	dictionary := newTrieNode()
	wordList := []string{"abc", "def", "ghi", "jkl"}
	for _, word := range wordList {
		dictionary.Insert(word)
	}
	textReader := strings.NewReader("Abc xxx DEF. \ngHi jkl zzz")

	spellingErrors := checkReaderConcurrent(dictionary, textReader)

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

//func TestCheckLineChan(t *testing.T) {
//
//}
