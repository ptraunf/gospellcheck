package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

func usage() {
	fmt.Println("\nUsage:\n\tgospellcheck WORDLIST [FILE | -]\n")
}

type SpellingError struct {
	misspelled   string
	line         int
	sentence     int
	wordPosition int
}

func (se SpellingError) String() string {
	return fmt.Sprintf("Line %d, sentence %d, word %d: '%s'", se.line, se.sentence, se.wordPosition, se.misspelled)
}

func normalizeWord(word string) string {
	re, err := regexp.Compile("[^a-z]")
	if err != nil {
		log.Fatal(err)
	}
	normalizedBytes := re.ReplaceAll([]byte(strings.ToLower(word)), []byte(""))
	return string(normalizedBytes)
}

func checkLine(dictionary *Trie, line string, lineNum int) []SpellingError {
	spellingErrors := make([]SpellingError, 0)
	sentences := strings.FieldsFunc(line, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})
	for s, sentence := range sentences {
		trimmedSentence := strings.Trim(sentence, ". ")
		if len(trimmedSentence) == 0 {
			continue
		}
		words := strings.Split(trimmedSentence, " ")
		for w, word := range words {
			normalized := normalizeWord(word)
			if len(normalized) > 0 && !dictionary.Contains(normalized) {
				spellingErrors = append(spellingErrors, SpellingError{
					misspelled:   word,
					line:         lineNum,
					sentence:     s + 1,
					wordPosition: w + 1,
				})
			}
		}
	}
	return spellingErrors
}

func checkLineChan(dictionary *Trie, line string, lineNum int, out chan<- SpellingError) {
	sentences := strings.FieldsFunc(line, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})
	for s, sentence := range sentences {

		trimmedSentence := strings.Trim(sentence, ". ")
		if len(trimmedSentence) == 0 {
			continue
		}
		words := strings.Split(trimmedSentence, " ")

		for w, word := range words {
			normalized := normalizeWord(word)
			if len(normalized) > 0 && !dictionary.Contains(normalized) {
				out <- SpellingError{
					misspelled:   word,
					line:         lineNum,
					sentence:     s + 1,
					wordPosition: w + 1,
				}
			}
		}
	}
}

func checkChannel(dictionary *Trie, lines <-chan string) chan SpellingError {
	errChan := make(chan SpellingError)

	go func() {
		defer close(errChan)
		i := 0
		for line := range lines {

			checkLineChan(dictionary, line, i+1, errChan)
			i++

		}
	}()

	return errChan
}
func checkReaderChannel(dictionary *Trie, reader io.Reader) chan SpellingError {
	linesChan := make(chan string)
	scanner := bufio.NewScanner(reader)
	go func() {

		defer close(linesChan)
		for scanner.Scan() {
			line := scanner.Text()
			linesChan <- line
		}
	}()

	errChan := checkChannel(dictionary, linesChan)
	return errChan
}
func checkReader(dictionary *Trie, reader io.Reader) []SpellingError {
	spellingErrors := make([]SpellingError, 0)

	scanner := bufio.NewScanner(reader)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineErrors := checkLine(dictionary, line, i+1)
		spellingErrors = append(spellingErrors, lineErrors...)
		i++
	}
	slices.SortFunc(spellingErrors, func(a, b SpellingError) int {
		return a.line - b.line
	})
	return spellingErrors
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

	dict, err := initializeDictionary(f)
	if err != nil {
		log.Fatal(err)
	}
	if "-" == os.Args[2] {

		in := bufio.NewReader(os.Stdin)
		spellingErrors := checkReader(dict, in)
		for _, spellingError := range spellingErrors {
			fmt.Printf("%s\n", spellingError.String())
		}
		return
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

		fileReader := bufio.NewReader(targetFile)

		spellingErrors := checkReaderChannel(dict, fileReader)
		//defer close(spellingErrors)
		for spellingError := range spellingErrors {
			fmt.Printf("%s\n", spellingError.String())
		}

	}
}
