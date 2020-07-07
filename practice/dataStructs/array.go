package main

import "fmt"

type myArray []string

func main() {
	newArray := myArray{}
	newArray = newArray.push("a", "b", "c")
	newArray = newArray.pop()
	newArray = newArray.push("d", "e")
	newArray = newArray.delete(1)
	fmt.Println(newArray.lookup(0))
	newArray = newArray.unshift("x", "y", "z")
	newArray = newArray.splice(1, "F", "G")
	fmt.Println(newArray)
}

//0(n)
func (a myArray) push(data ...string) myArray {
	return append(a, data...)
}

//0(1)
func (a myArray) pop() myArray {
	if len(a) < 1 {
		return nil
	}

	return a[:len(a)-1]
}

//0(n)
func (a myArray) delete(index int) myArray {
	return append(a[:index], a[index+1:]...)
}

//0(1)
func (a myArray) lookup(index int) string {
	return a[index]
}

//0(n)
//prepend
func (a myArray) unshift(data ...string) myArray {
	return append(data, a...)
}

//0(n)
//put data at start index
func (a myArray) splice(start int, data ...string) myArray {
	copy(a[start:], data)
	return a
}
