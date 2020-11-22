package main

import (
	"log"

	"github.com/lets_Go/practice/todo/cmd"
	"github.com/lets_Go/practice/todo/db"
)

func main() {
	if err := db.NewDB(); err != nil {
		log.Fatal(err)
	}

	cmd.RootCmd.Execute()
}
