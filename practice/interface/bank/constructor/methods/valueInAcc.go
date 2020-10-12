package methods

import "fmt"

func Value(p Payer) {
	switch p.(type) {
	case *Wallet:
		fmt.Println("amount left - ", p.(*Wallet).Amount)

	case *ApplePay:
		fmt.Println("amount left - ", p.(*ApplePay).Amount)

	case *Card:
		fmt.Println("amount left - ", p.(*Card).Amount)

	case *YandexMoney:
		fmt.Println("amount left - ", p.(*YandexMoney).Amount)
	}
}
