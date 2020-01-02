package main

import "log"

import "fmt"

type account struct {
	balance float64
}

func (a *account) Balance() float64 {
	return a.balance
}

func (a *account) deposit(amount float64) {
	log.Printf("Depositing: %f", amount)
	a.balance += amount
}

func (a *account) withdraw(amount float64) {
	if amount > a.balance {
		return
	}
	log.Printf("withdrawing: %f", amount)
	a.balance -= amount
}

func main() {
	var acc account
	var cash float64 = 50
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				if j%2 == 1 {
					acc.withdraw(cash)
					fmt.Printf("customer %d withdrew %f\n", j, cash)
					continue
				}
				acc.deposit(50)
				fmt.Printf("customer %d deposited %f\n", j, cash)
			}
		}()
	}
	fmt.Scanln()
	fmt.Println("\t\taccount balance\t\t",acc.Balance())
}
