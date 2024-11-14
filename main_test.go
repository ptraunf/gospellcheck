package main

import (
	"strings"
	"testing"
)

func TestValidateFilename_Valid(t *testing.T) {
	filename := "./main.go"
	_, err := validateFilename(filename)
	if err != nil {
		t.Fatalf("Expected no error; got:\n%v", err.Error())
	}
}

func TestValidateFilename_Invalid(t *testing.T) {
	filename := "../some/path/to/words.file"
	_, err := validateFilename(filename)
	if err == nil {
		t.Fatal("Expected error; got no error")
	}
}
func TestValidateFilename_DoesNotExist(t *testing.T) {
	filename := "./some/path/to/words.file"
	_, err := validateFilename(filename)
	if err == nil {
		t.Fatal("Expected error; got no error")
	}
}
func TestCheckLine(t *testing.T) {
	dictionary := newTrie()
	wordList := []string{"abc", "def", "ghi", "jkl"}
	for _, word := range wordList {
		dictionary.Insert(word)
	}
	text := "Abc xxx def yyy. Ghi zzz"
	spellingErrors := checkLine(dictionary, text, 1)

	if len(spellingErrors) < 3 {
		t.Fatalf("\nExpected 3 spelling errors; got %v\n", len(spellingErrors))
	}
	for _, spellingError := range spellingErrors {
		t.Logf("\n%s\n", spellingError)
	}
}

func TestCheckReader(t *testing.T) {
	dictionary := newTrie()
	wordList := []string{"abc", "def", "ghi", "jkl"}
	for _, word := range wordList {
		dictionary.Insert(word)
	}
	textReader := strings.NewReader("Abc xxx DEF. \ngHi jkl zzz")

	spellingErrors := checkReader(dictionary, textReader)

	if len(spellingErrors) != 2 {

		t.Fatalf("\nExpected 2 spelling errors; got %v\n", len(spellingErrors))
	}
	for _, spellingError := range spellingErrors {
		t.Logf("\n%v\n", spellingError)
	}

}
