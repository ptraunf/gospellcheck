package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"slices"
	"strings"
	"sync"
)

type SpellingError struct {
	misspelled   string
	line         int
	sentence     int
	wordPosition int
}
type Spellcheck interface {
	InitializeWordList(r io.Reader)
	CheckReader(r io.Reader) chan SpellingError
}

type TrieSpellcheck struct {
	trie Trie
}

func newSpellcheck() Spellcheck {
	return &TrieSpellcheck{}
}

func (spellcheck *TrieSpellcheck) InitializeWordList(r io.Reader) {
	spellcheck.trie = newTrieNode()
	spellcheck.trie.InsertAll(r)
}

func (spellcheck *TrieSpellcheck) CheckReader(r io.Reader) chan SpellingError {
	return checkReaderConcurrent(spellcheck.trie, r)
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

func checkLine(dictionary Trie, line string, lineNum int) []SpellingError {
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

func checkLineChan(dictionary Trie, line string, lineNum int, out chan<- SpellingError, wg *sync.WaitGroup) {
	sentences := strings.FieldsFunc(line, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})
	var sentenceWg sync.WaitGroup
	for s, sentence := range sentences {
		sentenceWg.Add(1)
		go func(sentenceNum int, sentence string) {
			defer sentenceWg.Done()
			//log.Printf("Reading sentence %d from line %d\n", s+1, lineNum)
			trimmedSentence := strings.Trim(sentence, ". ")
			if len(trimmedSentence) == 0 {
				return
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
					//log.Printf("Sentence %d yielded spelling error\n", sentenceNum)
				}
			}
		}(s, sentence)
	}
	sentenceWg.Wait()
	wg.Done()
}

func checkChannel(dictionary Trie, lines <-chan string) chan SpellingError {
	errChan := make(chan SpellingError)
	var wg sync.WaitGroup
	i := 0
	for line := range lines {
		wg.Add(1)
		//log.Printf("Read line %d from linesChan\n", i+1)
		go checkLineChan(dictionary, line, i+1, errChan, &wg)
		i++
	}
	go func() {
		defer close(errChan)
		wg.Wait()
	}()
	return errChan
}

func checkReaderConcurrent(dictionary Trie, reader io.Reader) chan SpellingError {
	linesChan := make(chan string)
	scanner := bufio.NewScanner(reader)
	go func() {
		defer close(linesChan)
		for scanner.Scan() {
			line := scanner.Text()
			//log.Printf("scanned line\n")
			linesChan <- line
		}
	}()

	errChan := checkChannel(dictionary, linesChan)
	return errChan
}
func checkReaderSequential(dictionary Trie, reader io.Reader) []SpellingError {
	//log.Printf("Spellcheck Sequential\n")

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
