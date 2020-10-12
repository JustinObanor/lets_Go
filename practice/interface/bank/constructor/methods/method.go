package methods

import "errors"

type Payer interface {
	Pay(amount int) error
}

type Wallet struct {
	Amount int
}

type Card struct {
	Name   string
	Amount int
}

type ApplePay struct {
	ID     string
	Amount int
}

type YandexMoney struct {
	ID string
	Card
}

func (a *ApplePay) Pay(amount int) error {
	if amount > a.Amount {
		return errors.New("insufficient funds")
	}
	a.Amount -= amount
	return nil
}

func (w *Wallet) Pay(amount int) error {
	if amount > w.Amount {
		return errors.New("insufficient funds")
	}
	w.Amount -= amount
	return nil
}

func (c *Card) Pay(amount int) error {
	if amount > c.Amount {
		return errors.New("insufficient funds")
	}
	c.Amount -= amount
	return nil
}

