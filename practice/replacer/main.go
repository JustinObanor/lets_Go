package main

import (
	"fmt"
	"strings"
)

func replaceQuotes(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r == '\'' {
			sb.WriteRune('"')
		} else if r == '"' {
			sb.WriteRune('\'')
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func main() {
	s := ` a : " ", b : ' ' `
	fmt.Println(replaceQuotes(s))

}
