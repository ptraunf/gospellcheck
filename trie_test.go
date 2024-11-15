package main

import (
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewTrieNode(t *testing.T) {
	var sc Trie = newTrieNode()
	if sc == nil {
		t.Fatalf("Could not create new Trie")
	}
}
func TestInsert(t *testing.T) {
	strs := []string{"ant", "anthem", "anteater"}
	var trie Trie = newTrieNode()
	for _, word := range strs {
		insertResult := trie.Insert(word)
		if !insertResult {
			t.Fatalf("\nExpected:\ttrue\nActual:\t%v", insertResult)
		}
	}

}
func TestEnumerate(t *testing.T) {
	strs := []string{"ant", "anthem", "anteater"}
	var trie *TrieNode = newTrieNode()
	var insertResult bool = true
	for _, word := range strs {
		insertResult = insertResult && trie.Insert(word)
	}
	results := trie.Enumerate()
	expectedNResults := len(strs)
	actualNResults := len(results)
	if expectedNResults != actualNResults {
		t.Fatalf("\nExpected %v results; got %v results", expectedNResults, actualNResults)
	}
	t.Logf("Trie:\n%v\n", trie)
}
func TestRemove(t *testing.T) {

}
func TestContains_Positive(t *testing.T) {
	strs := []string{"ant", "anthem", "anteater"}
	var trie Trie = newTrieNode()
	var insertResult bool = true
	for _, word := range strs {
		insertResult = insertResult && trie.Insert(word)
	}
	key := "ant"
	containsKey := trie.Contains(key)
	if !containsKey {
		t.Fatalf("Expected trie to contain %v\n", key)
	}
}
func TestContains_Negative(t *testing.T) {
	strs := []string{"ant", "anthem", "anteater"}
	var trie Trie = newTrieNode()
	var insertResult bool = true
	for _, word := range strs {
		insertResult = insertResult && trie.Insert(word)
	}
	key := "blarg"
	containsKey := trie.Contains(key)
	if containsKey {
		t.Fatalf("Did not expect trie to contain key %v\n", key)
	}
}
func TestLongestPrefix(t *testing.T) {
	s := "abcdefg"
	trie := newTrieNode()
	wordList := []string{"a", "ab", "abc", "abcd", "abcdex"}
	trie.InsertAll(strings.NewReader(strings.Join(wordList, "\n")))
	expected := "abcde"
	actual := trie.LongestPrefix(s)
	if actual != expected {
		t.Fatalf("\nExpected:\t%v\nActual:\t\t%v\n", expected, actual)
	}
}

func TestLongestPrefixNone(t *testing.T) {
	s := "abcdefg"
	trie := newTrieNode()
	wordList := []string{"x", "xy", "xyz"}
	trie.InsertAll(strings.NewReader(strings.Join(wordList, "\n")))
	expected := ""
	actual := trie.LongestPrefix(s)
	if actual != expected {
		t.Fatalf("\nExpected:\t%v\nActual:\t\t%v\n", expected, actual)
	}
}

func TestKeysWithPrefix(t *testing.T) {
	s := "abcdefg"
	trie := newTrieNode()
	wordList := []string{"a", "abx", "abcx", "abcdx", "abcdex", "abcdey", "abcdez"}
	trie.InsertAll(strings.NewReader(strings.Join(wordList, "\n")))
	expected := []string{"abcdex", "abcdey", "abcdez"}
	actual := trie.KeysWithCommonPrefix(s)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nExpected:\t%v\nActual:\t\t%v\n", expected, actual)
	}
}

func pickRandomWords(n int, t *TrieNode) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sample := make([]string, n)
	words := t.Enumerate()
	nWords := len(words)
	for i := 0; i < n; i++ {
		wordIdx := r.Intn(nWords)
		sample[i] = words[wordIdx]
	}
	return sample
}
func BenchmarkContains(b *testing.B) {
	f, err := os.Open("words.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	trie := newTrieNode()
	trie.InsertAll(f)
	lookupWords := pickRandomWords(10, trie)
	l := len(lookupWords)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := lookupWords[i%l]
		trie.Contains(word)
	}
}
func BenchmarkInsert(b *testing.B) {
	f, err := os.Open("words.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	trie1 := newTrieNode()
	trie1.InsertAll(f)
	list := trie1.Enumerate()
	l := len(list)
	trie2 := newTrieNode()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie2.Insert(list[i%l])
	}
}

func TestInsertAll(t *testing.T) {
	values := []string{"a", "b", "c", "d"}
	trie := newTrieNode()
	trie.InsertAll(strings.NewReader(strings.Join(values, "\n")))
	for _, v := range values {
		if !trie.Contains(v) {
			t.Fatalf("Trie should contain %v\n", v)
		}
	}
}
