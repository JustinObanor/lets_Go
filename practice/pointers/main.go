package main

func main() {
	a := 2
	b := &a 
	*b = 3 
	c := &a

	d := new(int)
	*d = 12
	*c = *d  
	*d = 13 

	c = d
	*c = 14
}
