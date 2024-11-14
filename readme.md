# Spellcheck

## Usage:
```sh
gospellcheck WORDLIST [FILE|-]
```

## Examples
### Spellcheck a file:
```sh
gospellcheck words.txt my_content.txt
# Outputs:
Line 1, sentence 1, word 1: 'Thes'
Line 4, sentence 1, word 2: 'wrds'
Line 4, sentence 1, word 6: 'inkorrect'
```
### Spellcheck text read from stdin
```sh
pdftotext my_content.pdf - | gospellcheck words.txt -
# Outputs:
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

