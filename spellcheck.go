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
	suggestions  []string
}
type Spellcheck interface {
	InitializeWordList(r io.Reader)
	CheckReader(r io.Reader) chan SpellingError
	GetSuggestions(word string) []string
}

type TrieSpellcheck struct {
	trie         Trie
	nSuggestions int
}

func newSpellcheck(nSuggestions int) Spellcheck {

	return &TrieSpellcheck{
		nSuggestions: nSuggestions,
	}
}

func (spellcheck *TrieSpellcheck) InitializeWordList(r io.Reader) {
	spellcheck.trie = newTrieNode()
	spellcheck.trie.InsertAll(r)
}

func (spellcheck *TrieSpellcheck) CheckReader(r io.Reader) chan SpellingError {
	linesChan := make(chan string)
	scanner := bufio.NewScanner(r)
	go func() {
		defer close(linesChan)
		for scanner.Scan() {
			line := scanner.Text()
			linesChan <- line
		}
	}()

	errChan := spellcheck.checkLines(linesChan)
	return errChan
}

func (se SpellingError) String() string {
	s := fmt.Sprintf("Line %d, sentence %d, word %d: '%s'", se.line, se.sentence, se.wordPosition, se.misspelled)
	if len(se.suggestions) > 0 {
		s = s + fmt.Sprintf("\n\tSuggestions: %v", se.suggestions)
	}
	return s
}

func normalizeWord(word string) string {
	re, err := regexp.Compile("[^a-z]")
	if err != nil {
		log.Fatal(err)
	}
	normalizedBytes := re.ReplaceAll([]byte(strings.ToLower(word)), []byte(""))
	return string(normalizedBytes)
}
func (spellcheck *TrieSpellcheck) GetSuggestions(word string) []string {
	suggestions := spellcheck.trie.KeysWithCommonPrefix(word)
	if len(suggestions) > spellcheck.nSuggestions {
		return suggestions[:spellcheck.nSuggestions]
	} else {
		return suggestions
	}
}

func (spellcheck *TrieSpellcheck) checkLine(line string, lineNum int, out chan<- SpellingError, wg *sync.WaitGroup) {
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
				if len(normalized) > 0 && !spellcheck.trie.Contains(normalized) {
					spellingError := SpellingError{
						misspelled:   word,
						line:         lineNum,
						sentence:     sentenceNum + 1,
						wordPosition: w + 1,
					}
					if spellcheck.nSuggestions > 0 {
						spellingError.suggestions = spellcheck.GetSuggestions(normalized)
					}
					out <- spellingError
				}
			}
		}(s, sentence)
	}
	sentenceWg.Wait()
	wg.Done()
}

func (spellcheck *TrieSpellcheck) checkLines(lines <-chan string) chan SpellingError {
	errChan := make(chan SpellingError)
	var wg sync.WaitGroup
	i := 0
	for line := range lines {
		wg.Add(1)
		go spellcheck.checkLine(line, i+1, errChan, &wg)
		i++
	}
	go func() {
		defer close(errChan)
		wg.Wait()
	}()
	return errChan
}
