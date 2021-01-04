package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var re = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]) of ([0-9]{3})[)][.](.+?)$")

var dir = "sample"
var dry = flag.Bool("dry", true, "specify whether its a dry run or not")

func main() {
	if err := renameFiles(dir); err != nil {
		log.Panic(err)
	}
}

func renameFiles(root string) error {
	var replaceStrA = "$2 - $1 - $3 of $4.$5"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if match(info.Name()) == "" {
			return nil
		}

		path = filepath.Dir(path)
		oldpath := filepath.Join(path, info.Name())
		newpath := filepath.Join(path, re.ReplaceAllString(info.Name(), replaceStrA))

		fmt.Printf("mv %s -> %s\n", oldpath, newpath)

		if !*dry {
			if err := os.Rename(oldpath, newpath); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func match(filename string) string {
	return re.FindString(filename)
}
