package main

import (
	"testing"
)

func TestMain(t *testing.T) {

}

func TestValidateFilename_Valid(t *testing.T) {
	filename := "some/path/to/words.file"
	err := sanitizeFilename(filename)
	if err != nil {
		t.Fatalf("Expected no error; got:\n%v", err.Error())
	}
}

func TestValidateFilename_Invalid(t *testing.T) {
	filename := "../some/path/to/words.file"
	err := sanitizeFilename(filename)
	if err == nil {
		t.Fatal("Expected error; got no error")
	}
}
