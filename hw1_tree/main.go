package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	fileToIgnore = ".DS_Store"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// printFiles = false
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	// fmt.Println(directory)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()
		if fileName == fileToIgnore {
			continue
		}

		if printFiles || fileInfo.IsDir() {
			if fileInfo.Mode().IsDir() {
				fmt.Printf("├───%s\n", fileName)
			} else {
				fmt.Println(fileName)
			}

			if err := os.Chdir(path); err != nil {
				return err
			}

			dirTree(out, fileName, printFiles)

			if err := os.Chdir("../"); err != nil {
				return err
			}
		}
	}

	return nil
}
