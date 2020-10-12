package constructor

import (
	"github.com/lets_Go/test/methods"
)

func NewCard(name string, amount int) *methods.Card {
	return &methods.Card{
		Name:   name,
		Amount: amount,
	}
}

func NewWallet(name string, amount int) *methods.Wallet {
	return &methods.Wallet{
		Amount: amount,
	}
}

func NewApplePay(id string, amount int) *methods.ApplePay {
	return &methods.ApplePay{
		ID:     id,
		Amount: amount,
	}
}

func NewYandexMoney(id, name string, amount int) *methods.YandexMoney {
	return &methods.YandexMoney{
		ID: id,
		Card: methods.Card{
			Name:   name,
			Amount: amount,
		},
	}
}
