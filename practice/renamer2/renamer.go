package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var reA = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]) of ([0-9]{3})[)][.](.+?)$")
var replaceStrA = "$2 - $1 - $3 of $4.$5"

var reB = regexp.MustCompile("^(.+?)_([0-9]{3})[.](.+?)$")
var replaceStrB = "$1 - $2 of 100.$3"

var dir = "sample"
var dry = flag.Bool("dry", true, "specify whether its a dry run or not")

func main() {
	flag.Parse()

	if err := renameFiles(dir); err != nil {
		log.Panic(err)
	}
}

func renameFiles(root string) error {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		newfileName := match(info.Name())
		if newfileName == "" {
			return nil
		}

		path = filepath.Dir(path)
		oldpath := filepath.Join(path, info.Name())
		newpath := filepath.Join(path, newfileName)

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
	switch {
	case reA.MatchString(filename):
		return reA.ReplaceAllString(filename, replaceStrA)
	case reB.MatchString(filename):
		return reB.ReplaceAllString(filename, replaceStrB)
	}
	return ""
}
