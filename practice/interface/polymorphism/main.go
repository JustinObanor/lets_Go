package main

import "fmt"

type Carer interface {
	detail() string
}

type car struct {
	name, brand string
	price       float64
}

type phone struct {
	model string
}

func (c car) detail() string {
	return fmt.Sprintf("%s of brand %s and a price tag of %.2f", c.name, c.brand, c.price)
}

func (p phone) detail() string {
	return fmt.Sprintf(" brand %s ", p.model)
}

func (c *car) changePrice(price float64) {
	c.price = price
}

func newCar(name, brand string, price float64) *car {
	return &car{
		name:  name,
		brand: brand,
		price: price,
	}
}

func print(c Carer) {
	// fmt.Println(c.detail())
}

func main() {
	//dynamic type - type checking is done at runtime
	//static type - type checking is done at compile time
	//interfaces have dynamic type that changes during runtime
	var cs Carer //abstract type value. nil type n value
	fmt.Printf("%T\n", cs)

	c := car{"ad", "sdf", 99.9} //concrete value of type car
	cs = c                      // interfac values are a pair of a concrete value and a dynamic type
	// its value holds a value of a specific underlying concrete type

	print(cs)
	//dynamic value is the value gotten at run time
	//dynamic type is the type gotten at run time
	//concrete type is main.car
	//polymorphism - cs can take many dynamic forms at runtime
	//dynamic type of cs can vary during execution
	fmt.Printf("%T\n", cs)

	p := phone{"Xiami"}
	cs = p
	fmt.Printf("%T\n", cs)
	//concrete type is main.phone. type assigned at runtime

	//interfaces isolate and decople parts of the program, like the exaple below
	// cs.changePrice() //compile time error
	//the method of the dynamical underlying type can be accessed through type assertion
	//	cs.(*car).changePrice(12.424)

}
