package main

import "fmt"

// Caesar Cipher
// (description more or less taken from Wikipedia)
//
//	In cryptography, a Caesar cipher, also known as Caesar's cipher,
//	the shift cipher, Caesar's code or Caesar shift, is one of the
//	simplest and most widely known encryption techniques. It is a type
//	of substitution cipher in which each letter in the plaintext is
//	replaced by a letter some fixed number of positions down the
//	alphabet. For example, with a left shift of 3, D would be replaced
//	by A, E would become B, and so on. The method is named after Julius
//	Caesar, who used it in his private correspondence.

// cipher takes in the text to be ciphered along with the direction that
// is being taken; -1 means encoding, +1 means decoding.
func cipher(shift rune, text string, direction int) string {
	// shift -> number of letters to move to right or left
	// offset -> size of the alphabet, in this case the plain ASCII
	offset := rune(26)

	// string->rune conversion
	runes := []rune(text)

	for index, char := range runes {
		// Iterate over all runes, and perform substitution
		// wherever possible. If the letter is not in the range
		// [1 .. 25], the offset defined above is added or
		// subtracted.
		switch direction {
		case -1: // decoding
			if char >= 'a'+shift && char <= 'z' || char >= 'A'+shift && char <= 'Z' {
				char = char - shift
			} else if char >= 'a' && char <= 'a'+shift || char >= 'A' && char <= 'A'+shift {
				char = char - shift + offset
			}
		case +1: // encoding
			if char >= 'a' && char <= 'z'-shift || char >= 'A' && char <= 'Z'-shift {
				char = char + shift
			} else if char >= 'z'-shift && char <= 'z' || char >= 'Z'-shift && char <= 'Z' {
				char = char + shift - offset
			}
		}

		// Above `if`s handle both upper and lower case ASCII
		// characters; anything else is returned as is (includes
		// numbers, punctuation and space).
		runes[index] = char
	}

	return string(runes)
}

// encode and decode provide the API for encoding and decoding text using
// the Caesar Cipher algorithm.
func encode(shift rune, text string) string { return cipher(shift, text, -1) }

func decode(shift rune, text string) string { return cipher(shift, text, +1) }

// A simple test
func main() {
	print("Input first name: ")
	var fName string
	fmt.Scan(&fName)

	print("Input last name: ")
	var lName string
	fmt.Scan(&lName)

	fmt.Print("Input key: ")
	var key rune

	text := fName + " " + lName

	fmt.Scan(&key)
	encoded := encode(key, text)
	println("Encoded: " + encoded)
	decoded := decode(key, encoded)
	println("Decoded: " + decoded)
}

//title list
//language
//platform
//program code
//how code works(where we put data and how we receive info)
//encode and decode
