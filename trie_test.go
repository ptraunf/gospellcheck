package main

import (
	"math/rand"
	"os"
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
	f, err := os.Open("words.txt")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	trie, err := initializeDictionary(f)
	if err != nil {
		b.Fatal(err)
	}
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
	defer f.Close()
	trie1, err := initializeDictionary(f)
	if err != nil {
		b.Fatal(err)
	}
	list := trie1.Enumerate()
	l := len(list)
	trie2 := newTrie()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie2.Insert(list[i%l])
	}
}

type mockReader struct {
	data []byte
}

func (mr mockReader) Read(buf []byte) (n int, err error) {
	l := len(mr.data)
	if l > len(buf) {
		n = copy(buf, mr.data[0:l])
	} else {
		n = copy(buf, mr.data)
	}
	return
}
func TestInitializeDictionary(t *testing.T) {
	r := mockReader{
		data: []byte("A"),
	}
	dict, err := initializeDictionary(r)
	if dict == nil {
		t.Fatalf("nil dict")
	}
	if err != nil {
		t.Fatalf(err.Error())
	}
}
