package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir := "./sample"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var toRename []string
	var toVisit []string
	var count int

	for _, file := range files {
		if file.IsDir() {
			toVisit = append(toVisit, filepath.Join(dir, file.Name()))
		} else {
			_, err := match(file.Name(), 0)
			if err == nil {
				count++

				toRename = append(toRename, file.Name())
			}
		}
	}

	for _, oldpath := range toRename {
		newpath, err := match(oldpath, count)
		if err != nil {
			panic(err)
		}

		fmt.Printf("mv %s => %s\n", oldpath, filepath.Join(dir, newpath))

		err = os.Rename(newpath, filepath.Join(dir, oldpath))
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(toVisit)
}

func match(filename string, max int) (string, error) {
	pieces := strings.Split(filename, ".")

	ext := pieces[len(pieces)-1]

	file := strings.Join(pieces[:len(pieces)-1], ".")
	pieces = strings.Split(file, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")

	num, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", errors.New("doesnt match")
	}

	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), num, max, ext), nil
}
