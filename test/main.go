package main

import (
	"fmt"
	"time"
)

var location *time.Location

func main() {
	now := time.Now().UTC()
	fmt.Println(now)
	//12:35:09.331984 +0000 UTC

	location, _ := time.LoadLocation("Europe/Minsk")

	now3 := time.Now().UTC().In(location).Format(time.RFC1123)
	fmt.Println(now3)
}
