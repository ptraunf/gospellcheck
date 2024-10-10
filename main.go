package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"regexp"

	// "log"
	"os"
	"strings"
)

func usage() {
	fmt.Print("Usage:\n\tgospellcheck [dictionary]")
}

func sanitizeFilename(filename string) (string, error) {
	// disallow `..` parent-directory traversal
	pathSep := os.PathSeparator
	pattern := fmt.Sprintf("\\.\\.\\%v", string(pathSep))
	fmt.Printf("Pattern: %v", pattern)
	match, err := regexp.Match(pattern, []byte(filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	if match {
		return filename, errors.New("Path contains '..'")
	}
	return filename, nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	wordFile := os.Args[1]
	sanitizedFilepath, err := sanitizeFilename(wordFile)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(sanitizedFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dict, err := initializeDictionary(f)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Spellcheck initialized with word list:\n%v\n", trie)
	for {
		fmt.Printf("Spellcheck word: ")
		in := bufio.NewReader(os.Stdin)
		s, _ := in.ReadString('\n')
		s = strings.TrimSpace(s)
		correct := dict.Contains(s)
		if correct {
			fmt.Printf("'%v' is in the dictionary.\n", s)
		} else {
			fmt.Printf("'%v' is NOT in the dictionary\n", s)
		}
	}
}
