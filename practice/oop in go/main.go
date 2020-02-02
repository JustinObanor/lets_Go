package main

import (
	"fmt"
)

const (
	WHITE = iota
	BLACK
	BLUE
	RED
	YELLOW
)

type Box struct {
	width, height, depth float64
	color                Color
}

type Color byte
type BoxList []Box //a slice of boxes

// method
func (b Box) Volume() float64 {
	return b.width * b.height * b.depth
}

// method with a pointer receiver
/*Because we want to change Box's color in this method.
Thus, if we don't use a pointer,
it will only change the value inside a copy of Box.*/
func (b *Box) SetColor(c Color) {
	b.color = c
}

// method
func (bl BoxList) BiggestsColor() Color {
	v := 0.00
	k := Color(WHITE)
	for _, b := range bl {
		//if b.Volume is > anything, assign b.Volume to v
		if b.Volume() > v {
			v = b.Volume()
			k = b.color
		}
	}
	return k
}

// method
func (bl BoxList) PaintItBlack() {
	for i, _ := range bl {
		bl[i].SetColor(BLACK)
	}
}

// method
//outputs the colors
func (c Color) String() string {
	strings := []string{"WHITE", "BLACK", "BLUE", "RED", "YELLOW"}
	return strings[c]
}

func main() {
	boxes := BoxList{
		Box{4, 4, 4, RED},
		Box{10, 10, 1, YELLOW},
		Box{1, 1, 20, BLACK},
		Box{10, 10, 1, BLUE},
		Box{10, 30, 1, WHITE},
		Box{20, 20, 20, YELLOW},
	}

	fmt.Printf("We have %d boxes in our set\n", len(boxes))
	fmt.Println("The volume of the first one is", boxes[0].Volume(), "cm³")
	fmt.Println("The color of the last one is", boxes[len(boxes)-1].color.String())
	fmt.Println("The biggest one is", boxes.BiggestsColor().String())

	// Let's paint them all black
	boxes.PaintItBlack()

	fmt.Println("The color of the second one is", boxes[1].color.String())
	fmt.Println("Obviously, now, the biggest one is", boxes.BiggestsColor().String())
}
