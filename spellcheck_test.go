package main

import (
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
