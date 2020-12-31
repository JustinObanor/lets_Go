package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var dir = "sample"

func main() {
	if err := rename(dir); err != nil {
		panic(err)
	}
}

func rename(directory string) error {
	toRename, err := getFiles(dir)
	if err != nil {
		return err
	}

	for _, files := range toRename {
		n := len(files)
		for idx, file := range files {

			oldfile := filepath.Join(file.path, file.file)
			res, _ := match(file.file)

			newfile := filepath.Join(file.path, fmt.Sprintf("%s %d of %d.%s", strings.Title(res.base), (idx+1), n, res.ext))

			fmt.Printf("mv %s -> %s\n", oldfile, newfile)
		}
	}

	return nil
}

type fileinfo struct {
	path string
	file string
}

func getFiles(root string) (map[string][]fileinfo, error) {
	toRename := make(map[string][]fileinfo)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		path = filepath.Dir(path)

		if file, err := match(info.Name()); err == nil {
			oldfile := fmt.Sprintf("%s_%s.%s", file.base, file.num, file.ext)
			key := filepath.Join(path, fmt.Sprintf("%s.%s", file.base, file.ext))

			toRename[key] = append(toRename[key], fileinfo{
				path: path,
				file: oldfile,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return toRename, nil
}

type matchResult struct {
	base string
	num  string
	ext  string
}

func match(filename string) (*matchResult, error) {
	pieces := strings.Split(filename, ".")
	file := strings.Join(pieces[:len(pieces)-1], ".")

	ext := pieces[len(pieces)-1]

	pieces = strings.Split(file, "_")

	file = strings.Join(pieces[:len(pieces)-1], "_")

	num := pieces[len(pieces)-1]

	_, err := strconv.Atoi(num)
	if err != nil {
		return nil, fmt.Errorf("%s doesnt match", filename)
	}

	return &matchResult{
		base: file,
		num:  num,
		ext:  ext,
	}, nil
}
