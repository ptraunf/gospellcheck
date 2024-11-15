package main

import (
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
