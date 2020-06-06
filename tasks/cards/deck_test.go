package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 16 {
		t.Errorf("Expected 16 got %d", len(d))
	}

	s := strings.Split(d[0], " ")[1]
	if s != "of" {
		t.Errorf("Expected of got %v", s)
	}
	fmt.Println(d[0])
}

func TestIOFile(t *testing.T) {
	os.Remove("_decktesting.txt")
	deck := newDeck()
	deck.saveToFile("_decktesting.txt")
	readDeck := readFromFile("_decktesting.txt")

	if len(readDeck) != 16 {
		t.Errorf("Expected 16 got %d", len(readDeck))
	}
	os.Remove("_decktesting.txt")
}
