package gospellcheck

import "testing"

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

func TestLongestPrefix(t *testing.T) {

}

func TestKeysWithPrefix(t *testing.T) {

}
