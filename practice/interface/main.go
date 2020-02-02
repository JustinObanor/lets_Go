package main

import "fmt"

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human
	school string
	loan   float32
}

type Employee struct {
	Human
	company string
	money   float32
}

type Men interface {
	sayHi()
}

func (h Human) sayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

func (e Employee) sayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone)
}

func (s Student) sayHi() {
	fmt.Printf("Hi, I am %s, I study at %s. I have $%f\n", s.name,
		s.school, s.loan)
}

func (s *Student) BorrowMoney(cash float32) {
	s.loan += cash
}

func (e *Employee) SpendCash(cash float32) {
	e.money -= cash
}

func main() {
	h1 := Human{"justin", 20, "123456789"}
	s1 := Student{Human{"james", 19, "12345678"}, "qwerty", 12.12}
	e1 := Employee{Human{"mike", 112, "546372883475"}, "ytrew", 12.345}
	h1.sayHi()
	s1.sayHi()
	e1.sayHi()
	s1.BorrowMoney(12.12)
	s1.sayHi()

	//var i Men
	// i = h1
	// i.sayHi()
	// i = s1
	// i.sayHi()
	// i = e1
	// i.sayHi()

	x := make([]Men, 3)
	x[0], x[1], x[2] = h1, e1, s1
	for _, v := range x {
		v.sayHi()
	}

}
