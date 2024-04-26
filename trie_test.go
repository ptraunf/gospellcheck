package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	var sc *Trie = newTrie()
	if sc == nil {
		t.Fatalf("StringContainer.New ")
	}
}
func TestInsert(t *testing.T) {
	strings := []string{"ant", "anthem", "anteater"}
	var trie *Trie = newTrie()
	for _, word := range strings {
		insertResult := trie.Insert(word)
		if !insertResult {
			t.Fatalf("\nExpected:\ttrue\nActual:\t%v", insertResult)
		}
	}

}
func TestEnumerate(t *testing.T) {
	strings := []string{"ant", "anthem", "anteater"}
	var trie *Trie = newTrie()
	var insertResult bool = true
	for _, word := range strings {
		insertResult = insertResult && trie.Insert(word)
	}
	results := trie.Enumerate()
	expectedNResults := len(strings)
	actualNResults := len(results)
	if expectedNResults != actualNResults {
		t.Fatalf("\nExpected %v results; got %v results", expectedNResults, actualNResults)
	}
	t.Logf("Trie:\n%v\n", trie)
}
func TestRemove(t *testing.T) {

}
func TestContains_Positive(t *testing.T) {
	strings := []string{"ant", "anthem", "anteater"}
	var trie *Trie = newTrie()
	var insertResult bool = true
	for _, word := range strings {
		insertResult = insertResult && trie.Insert(word)
	}
	key := "ant"
	containsKey := trie.Contains(key)
	if !containsKey {
		t.Fatalf("Expected trie to contain %v\n", key)
	}
}
func TestContains_Negative(t *testing.T) {
	strings := []string{"ant", "anthem", "anteater"}
	var trie *Trie = newTrie()
	var insertResult bool = true
	for _, word := range strings {
		insertResult = insertResult && trie.Insert(word)
	}
	key := "blarg"
	containsKey := trie.Contains(key)
	if containsKey {
		t.Fatalf("Did not expect trie to contain key %v\n", key)
	}
}

func testLongestPrefix(t *testing.T) {

}

func testKeysWithPrefix(t *testing.T) {

}

func pickRandomWords(n int, t *Trie) []string {
	rand.Seed(time.Now().UnixNano())
	sample := make([]string, n)
	words := t.Enumerate()
	nWords := len(words)
	for i := 0; i < n; i++ {
		wordIdx := rand.Intn(nWords)
		sample[i] = words[wordIdx]
	}
	return sample
}
func BenchmarkContains(b *testing.B) {
	trie := initializeSpellcheck("words.txt")
	lookupWords := pickRandomWords(10, trie)
	l := len(lookupWords)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := lookupWords[i%l]
		trie.Contains(word)
	}
}
func BenchmarkInsert(b *testing.B) {
	trie1 := initializeSpellcheck("words.txt")
	list := trie1.Enumerate()
	l := len(list)
	trie2 := newTrie()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		word := list[i%l]
		trie2.Insert(word)
	}
}
