package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func initializeSpellcheck(filename string) *Trie {
	var trie = newTrie()
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := scanner.Text()
		trie.Insert(word)
	}

	return trie
}
func main() {
	wordFile := os.Args[1]
	trie := initializeSpellcheck(wordFile)
	// fmt.Printf("Spellcheck initialized with word list:\n%v\n", trie)
	for {
		fmt.Printf("Spellcheck word: ")
		in := bufio.NewReader(os.Stdin)
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(s)
		correct := trie.Contains(s)
		if correct {
			fmt.Printf("'%v' is spelled correctly.\n", s)
		} else {
			fmt.Printf("'%v' is NOT spelled correctly.\n", s)
		}
	}

}
