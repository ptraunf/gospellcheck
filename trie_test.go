package gospellcheck

import "testing"

func TestCreate(t *testing.T) {
	var sc *Trie = newTrie()
	if sc == nil {
		t.Fatalf("StringContainer.New ")
	}
}
func TestInsert(t *testing.T) {
	// strings := []string{"ant", "anthem", "anteater"}
	testString := "and"
	var trie *Trie = newTrie()
	result := trie.Insert(testString)
	if !result {
		t.Fatalf("\nExpected:\ttrue\nActual:\t%v", result)
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
func TestContains(t *testing.T) {

}

func TestLongestPrefix(t *testing.T) {

}

func TestKeysWithPrefix(t *testing.T) {

}
