package main

func main() {
	cards := newDeck()

	cards.shuffle()

	cards.saveToFile("file.txt")

	readFromFile("file.txt").print()
}
