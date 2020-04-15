package main

import (
	"fmt"
	"strings"
)

type MyString string

type FindVowel interface {
	findVowel() []rune
	countVowel(s []rune)
}

func (m MyString) findVowel() []rune {
	var vowels []rune
	for _, vowel := range m {
		if vowel == 'a' || vowel == 'e' || vowel == 'i' || vowel == 'o' || vowel == 'u' {
			vowels = append(vowels, vowel)
		}
	}
	return vowels
}

func (m MyString) countVowel(s []rune) int {
	numOfa := strings.Count(string(s), "a")
	fmt.Printf(" \n\t has %d characters of [a]", numOfa)

	numOfe := strings.Count(string(s), "e")
	fmt.Printf("\n\t has %d characters of [e]", numOfe)

	numOfi := strings.Count(string(s), "i")
	fmt.Printf("\n\t has %d characters of [i]", numOfi)

	numOfo := strings.Count(string(s), "o")
	fmt.Printf("\n\t has %d characters of [o]", numOfo)

	numOfu := strings.Count(string(s), "u")
	fmt.Printf("\n\t has %d characters of [u]", numOfu)

	sum := numOfa + numOfe + numOfi + numOfo + numOfu
	return sum
}

func total(s ...int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	fmt.Println("\n Sum of vowels is", sum)
}

func main() {
	var s MyString
	s = "Hey! Lets see how many vowels we got "

	s.findVowel()

	v := s.findVowel()
	fmt.Printf("Vowels are %c", v)

	sum := s.countVowel(v)
	total(sum)
}
