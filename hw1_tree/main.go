package main

import (
	"fmt"
	"io"
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
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFile bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()

		if fileName == fileToIgnore {
			continue
		}

		if printFile {
			fmt.Println(fileName)

			if fileInfo.IsDir() {
				//wd - testdata
				if err := os.Chdir(path); err != nil {
					return err
				}
				dirTree(out, fileName, false)
				fmt.Println(fileName)

				if err := os.Chdir(fmt.Sprint("../")); err != nil {
					return err
				}
			}
		} 
	}
	return nil
}
