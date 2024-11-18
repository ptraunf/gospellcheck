package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func usage() {
	fmt.Printf("\nUsage:\n\tgospellcheck [OPTIONS] WORDLIST TARGET\n")
	fmt.Printf("\nWORDLIST\n\tnewline-delimited file of words to populate the spellcheck dictionary\n")
	fmt.Printf("\nTARGET\n\tfile to spellcheck, or '-' to read from stdin\n")
	fmt.Printf("\nOPTIONS\n\t-s\tnumber of words to suggest for each misspelling\n")
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

	suggestions := flag.Int("s", 0, "number of words to suggest for each misspelling")
	flag.Parse()
	if flag.NArg() < 2 {
		usage()
		return
	}
	wordFile := flag.Arg(0)
	targetPath := flag.Arg(1)

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

	spellcheck := newSpellcheck(*suggestions)
	spellcheck.InitializeWordList(f)

	var targetReader io.Reader
	if "-" == targetPath {

		targetReader = bufio.NewReader(os.Stdin)

	} else {
		//targetPath := os.Args[2]
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
