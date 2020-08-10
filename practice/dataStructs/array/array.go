package main

type myArray []string

//0(1)
func (a myArray) push(data string) myArray {
	return append(a, data)
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
