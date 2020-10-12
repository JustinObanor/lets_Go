package main

import (
	"errors"
	"fmt"
)

type Payer interface {
	Pay(n int) error
}

type Wallet struct {
	cash int
}

func (w *Wallet) Pay(n int) error {
	if n > w.cash {
		return errors.New("insuf env")
	}
	w.cash -= n
	return nil
}

func Buy(i interface{}) {
	var p Payer
	var ok bool

	if p, ok = i.(Payer); !ok {
		fmt.Println("cant pay with this", i)
		return
	}

	err := p.Pay(30)
	if err != nil {
		fmt.Println("error in payment", p, err)
		return
	}

	fmt.Println("thankjs for payinfg", p)
}

func main() {
	w := &Wallet{cash: 40}
	Buy(w)
	Buy([]int{1, 2, 3})
	Buy(2.1)
}
