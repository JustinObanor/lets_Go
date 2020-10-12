package methods

import "fmt"

func Buy(p Payer) {
	switch p.(type) {
	case *Wallet:
		wallet, ok := p.(*Wallet)
		if !ok {
			fmt.Println("couldnt get wallet")
		}
		fmt.Printf("using %T\n", wallet)

	case *Card:
		card, ok := p.(*Card)
		if !ok {
			fmt.Println("couldnt get card")
		}
		fmt.Printf("using %T\n", card)

	case *ApplePay:
		applePay, ok := p.(*ApplePay)
		if !ok {
			fmt.Println("couldnt get apple pay")
		}
		fmt.Printf("using %T\n", applePay)
	}

	err := p.Pay(20)
	if err != nil {
		fmt.Printf("Error paying in %T: %v\n\n", p, err)
		return
	}
	fmt.Printf("Thanks for buying through %T\n\n", p)
}
