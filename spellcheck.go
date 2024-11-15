package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
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
	return checkReader(spellcheck.trie, r)
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

func checkLine(trie Trie, line string, lineNum int, out chan<- SpellingError, wg *sync.WaitGroup) {
	sentences := strings.FieldsFunc(line, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})
	var sentenceWg sync.WaitGroup
	for s, sentence := range sentences {
		sentenceWg.Add(1)
		go func(sentenceNum int, sentence string) {
			defer sentenceWg.Done()
			trimmedSentence := strings.Trim(sentence, ". ")
			if len(trimmedSentence) == 0 {
				return
			}
			words := strings.Split(trimmedSentence, " ")

			for w, word := range words {
				normalized := normalizeWord(word)
				if len(normalized) > 0 && !trie.Contains(normalized) {
					out <- SpellingError{
						misspelled:   word,
						line:         lineNum,
						sentence:     sentenceNum + 1,
						wordPosition: w + 1,
					}
				}
			}
		}(s, sentence)
	}
	sentenceWg.Wait()
	wg.Done()
}

func checkLines(trie Trie, lines <-chan string) chan SpellingError {
	errChan := make(chan SpellingError)
	var wg sync.WaitGroup
	i := 0
	for line := range lines {
		wg.Add(1)
		go checkLine(trie, line, i+1, errChan, &wg)
		i++
	}
	go func() {
		defer close(errChan)
		wg.Wait()
	}()
	return errChan
}

func checkReader(trie Trie, reader io.Reader) chan SpellingError {
	linesChan := make(chan string)
	scanner := bufio.NewScanner(reader)
	go func() {
		defer close(linesChan)
		for scanner.Scan() {
			line := scanner.Text()
			linesChan <- line
		}
	}()

	errChan := checkLines(trie, linesChan)
	return errChan
}
