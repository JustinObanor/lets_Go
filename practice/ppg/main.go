package main

import (
	"errors"
	"fmt"
)

func makeGirls(ing int, ingReq, ingPres []int) (int, error) {
	if ing <= 1 || ing >= 1e7 {
		return 0, errors.New("violates constraing")
	}

	if len(ingReq) != ing || len(ingPres) != ing {
		return 0, errors.New("number of ingredients should be the same as quantity")
	}

	max := 0
	num := 0

	for {
		for i, v := range ingPres {
			num += ingReq[i]

			if num >= v {
				return max, nil
			}
			max++
		}
	}
}

func main() {
	fmt.Println("Input number of ingredients")
	var numOfIng int
	fmt.Scanf("%d", &numOfIng)

	fmt.Println("Input quantity of each ingredient required")
	ingReq := make([]int, numOfIng)
	for i := 0; i < numOfIng; i++ {
		fmt.Scanf("%d", &ingReq[i])
	}

	fmt.Println("Input quantity of each ingredient present")
	ingPres := make([]int, numOfIng)
	for i := 0; i < numOfIng; i++ {
		fmt.Scanf("%d", &ingPres[i])
	}

	if numOfIng <= 1 || numOfIng >= 1e7 {
		return
	}

	if len(ingReq) != numOfIng || len(ingPres) != numOfIng {
		return
	}

	max := 0
	num := 0

	for {
		for i, v := range ingPres {
			num += ingReq[i]

			if num >= v {
				fmt.Println(max)
			}
			max++
		}
	}

}
