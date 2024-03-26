# Spellcheck

Usage:
```sh
gospellcheck <word-list>
```

To generate wordlist on linux with `aspell` installed:
```sh
aspell dump master | tr 'A-Z' 'a-z' > words.txt
```
