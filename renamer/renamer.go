package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var dir = "sample"
var dry = flag.Bool("dry", true, "specify whether this should be a dry run or not")

func main() {
	flag.Parse()

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

			oldpath := filepath.Join(file.path, file.file)
			res, _ := match(file.file)

			newpath := filepath.Join(file.path, fmt.Sprintf("%s - %d of %d.%s", strings.Title(res.base), (idx+1), n, res.ext))

			fmt.Printf("mv %s -> %s\n", oldpath, newpath)
			if !*dry {
				if err := os.Rename(oldpath, newpath); err != nil {
					return err
				}
			}
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
			key := filepath.Join(path, fmt.Sprintf("%s.%s", file.base, file.ext))

			toRename[key] = append(toRename[key], fileinfo{
				path: path,
				file: info.Name(),
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
		ext:  ext,
	}, nil
}
