package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type file struct {
	name string
	path string
}

func main() {
	dir := "sample"

	err := rename(dir)
	if err != nil {
		panic(err)
	}
}

func rename(dir string) error {
	toRename, err := getFiles(dir)
	if err != nil {
		return err
	}

	for _, v := range toRename {
		newName, err := match(v.name, 0)
		if err != nil {
			return fmt.Errorf("error matching file %s: %s", v.name, err)
		}

		oldpath := filepath.Join(v.path, v.name)
		newpath := filepath.Join(v.path, newName)

		if err = os.Rename(oldpath, newpath); err != nil {
			return fmt.Errorf("error renaming file %s: %s", v.name, err)
		}
	}

	return nil
}

func getFiles(dir string) ([]file, error) {
	var f file
	var toRename []file

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			f.path = path
		}

		if _, err = match(info.Name(), 0); err == nil {
			f.name = info.Name()

			toRename = append(toRename, f)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return toRename, nil
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
