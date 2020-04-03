package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type deck []string

func newDeck() deck {
	cards := deck{}

	suits := []string{"spades", "diamonds", "heards", "clut"}
	values := []string{"ace", "two", "three", "four"}

	for _, s := range values {
		for _, v := range suits {
			s = strings.Title(s)
			v = strings.Title(v)

			cards = append(cards, s+" of "+v)
		}
	}
	return cards
}

func (d deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
}

// func (d deck) shuffle2(){
// 	s := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(s)
// 	randNum := r.Intn(len(d)-1)
// }

func (d deck) print() {
	for _, v := range d {
		fmt.Printf("%s\n", v)
	}
}

func deal(d deck, size int) (deck, deck) {
	return d[:size], d[size:]
}

func (d deck) toString() string {
	return strings.Join([]string(d), ",")
}

func readFromFile(filename string) deck {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	ss := strings.Split(string(bs), ",")
	return deck(ss)
}

func (d deck) saveToFile(filename string) error {
	if err := ioutil.WriteFile(filename, []byte(d.toString()), 0666); err != nil {
		return err
	}
	return nil
}
