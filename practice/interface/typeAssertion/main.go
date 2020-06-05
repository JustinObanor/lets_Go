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
	var c Carer = &car{"sda", "ads", 15.6}
	fmt.Printf("%T\n", c)

	// c.changePrice() cant access method

	//c.(*car).changePrice(324) type assertion

	// cc, ok := c.(car)//if true, cc holds the dynamic value
	// cc.changePrice(12314)
	// fmt.Println(cc, ok)

	//cases for c are compared based on the given interface value
	switch value := c.(type) {
	case *car:
		fmt.Println("type car")
		value.changePrice(999)
		fmt.Println(c.detail())
		// fmt.Println(value)
	}
}

/*

   // declaring an interface value that holds a circle type value
   var s shape = circle{radius: 2.5}

   fmt.Printf("%T\n", s) //interface dynamic type is circle

   // no direct access to interface's dynamic values
   // s.volume() -> error

   // there is access only to the methods that are defined inside the interface
   fmt.Printf("Circle Area:%v\n", s.area())

   // an interface value hides its dynamic value.
   // use type assertion to extract and return the dynamic value of the interface value.
   fmt.Printf("Sphere Volume:%v\n", s.(circle).volume())

   // checking if the assertion succeded or not
   ball, ok := s.(circle)
   if ok == true {
       fmt.Printf("Ball Volume:%v\n", ball.volume())
   }

    TYPE SWITCHES

   // it permits several type assertions in series
   switch value := s.(type) {
   case circle:
       fmt.Printf("%#v has circle type\n", value)
   case rectangle:
       fmt.Printf("%#v has rectangle type\n", value)

   }*/
