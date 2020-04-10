package main

<<<<<<< HEAD
import "fmt"

func main() {
	commits := map[int]int{
		1: 3711,
		2: 2138,
		3: 1908,
		4: 912,
	}
	for k, v := range commits {
		fmt.Println(k, v)
	}
=======
import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(11)
	for i := 0; i <= 10; i++ {
		
		go func(i int) {
			defer wg.Done()
			fmt.Printf("loop i is - %d\n", i)
		}(i)
	}
	wg.Wait()
	fmt.Println("Hello, Welcome to Go")
>>>>>>> 76da20ff0db08f457eb6e2096e51178dc80f05cf
}
