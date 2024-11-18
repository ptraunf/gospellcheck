# Spellcheck

## Usage:
```sh
gospellcheck [OPTIONS] WORDLIST TARGET
```
### Arguments
- `WORDLIST`: A file of words to populate the spellcheck dictionary, separated by new-lines
- `TARGET`: file to spellcheck, or '-' to read from stdin
- `OPTIONS`
  - `-s`: (integer) number of suggested words to include with each misspelling

## Examples
### Spellcheck a file:
```sh
gospellcheck words.txt my_content.txt
```
Outputs:
```
Line 1, sentence 1, word 1: 'Thes'
Line 4, sentence 1, word 2: 'wrds'
Line 4, sentence 1, word 6: 'inkorrect'
```

### Spellcheck with Suggestions:
```
gospellcheck -s 3 words.txt my_content.txt
```
Outputs 3 suggestions for each misspelling
```
Line 1, sentence 1, word 1: 'Thes'
	Suggestions: [thesauri thesaurus thesaurus's]
Line 1, sentence 1, word 2: 'wrds'
	Suggestions: [wrack wrack's wracked]
Line 1, sentence 1, word 4: 'inkorrect'
	Suggestions: [ink ink's inkblot]
```

### Spellcheck text read from stdin
```sh
pdftotext my_content.pdf - | gospellcheck words.txt -
```
Outputs
```
Line 1, sentence 1, word 1: 'Thes'
Line 4, sentence 1, word 2: 'wrds'
Line 4, sentence 1, word 6: 'inkorrect'
```

## Installation
```sh
git clone https://github.com/ptraunf/gospellcheck.git
cd gospellcheck
go install
```

### Generate a Word List
On linux/mac with `aspell` installed:
```sh
# gospellcheck is case-insensitive
aspell dump master | tr 'A-Z' 'a-z' | sort -u | uniq > words.txt
```

