package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func usage() {
	fmt.Println("\nUsage:\n\tgospellcheck WORDLIST [FILE | -]")
}

func validateFilename(filename string) (string, error) {
	pathSep := os.PathSeparator
	pattern := fmt.Sprintf("\\.\\.\\%v", string(pathSep))
	match, err := regexp.Match(pattern, []byte(filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	if match {
		return filename, errors.New("path contains '..'")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return filename, errors.New(fmt.Sprintf("file %s does not exist", filename))
	}
	return filename, nil
}

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}
	wordFile := os.Args[1]
	sanitizedFilepath, err := validateFilename(wordFile)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(sanitizedFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(f)

	spellcheck := newSpellcheck()
	spellcheck.InitializeWordList(f)

	var targetReader io.Reader
	if "-" == os.Args[2] {

		targetReader = bufio.NewReader(os.Stdin)

	} else {
		targetPath := os.Args[2]
		validTargetPath, err := validateFilename(targetPath)
		if err != nil {
			log.Fatal(err)
		}

		targetFile, err := os.Open(validTargetPath)
		defer func(f *os.File) {
			err := targetFile.Close()
			if err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}(targetFile)

		targetReader = bufio.NewReader(targetFile)

	}
	spellingErrorsCh := spellcheck.CheckReader(targetReader)
	spellingErrorsSlice := chanToSortedSlice(spellingErrorsCh, func(a, b SpellingError) int {
		return a.line - b.line
	})
	for _, spellingError := range spellingErrorsSlice {
		fmt.Printf("%s\n", spellingError.String())
	}
}
