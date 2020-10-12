package main

import (
	"github.com/lets_Go/test/constructor"
	"github.com/lets_Go/test/methods"
)

func main() {
	apple := constructor.NewApplePay("3034102u130u452", 19)
	// card := constructor.NewCard("Justin Obanor", 40)
	// wallet := constructor.NewWallet("Justin Obanor", 50)
	// yandex := constructor.NewYandexMoney("4546123435465524", "Justin Obanor", 30)
	
	methods.Value(apple)

	methods.Buy(apple)

	methods.Value(apple)
}
