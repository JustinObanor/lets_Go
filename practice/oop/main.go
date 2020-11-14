package main

import "fmt"

//abstraction is hiding tjhe internal implementations and details, and only showing whats neccessary.
//eg a house. hiding whats going on inside and only showing the neccasary things

//encapsulation(grouping data and methods together)
type person struct {
	fname, lname string
	age          int
}

//inheritance(embedding)
type student struct {
	person
}

//teacher inherits interface
type teacher struct {
	person
	salary int
}

func (p *person) setName(fname, lname string) {
	p.fname = fname
	p.lname = lname
}

type setNamer interface {
	setName(fname, lname string)
}

//polymorphism. interface supports polymorphism. different behavior based on the underlying type
func polyMorpher(i setNamer, fname, lname string) {
	switch i.(type) {
	case *teacher:
		fmt.Printf("%T", i)
	case *person:
		fmt.Printf("%T", i)
	case *student:
		fmt.Printf("%T", i)
	}

	i.setName(fname, lname)
}

func main() {
	p1 := &person{}

	s1 := &student{}

	t1 := &teacher{}

	polyMorpher(p1, "james", "mak")
	fmt.Println(p1)
	polyMorpher(s1, "john", "walt")
	fmt.Println(s1)
	polyMorpher(t1, "mike", "mikson")
	fmt.Println(t1)

}
