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

func rename(dir string) error {
	toRename, err := getFiles(dir)
	if err != nil {
		return err
	}

	for _, files := range toRename {
		n := len(files)
		for idx, file := range files {
			oldpath := filepath.Join(file.path, fmt.Sprintf("%s.%s", file.base, file.ext))
			newpath := filepath.Join(file.path, fmt.Sprintf("%s %d of %d.%s", strings.Title(file.base), (idx+1), n, file.ext))

			if err := os.Rename(oldpath, newpath); err != nil {
				return fmt.Errorf("error renaming file %s to %s: %v", oldpath, newpath, err)
			}
		}
	}
	return nil
}

type fileinfo struct {
	path string
	base string
	ext  string
}

func getFiles(root string) (map[string][]fileinfo, error) {
	toRename := make(map[string][]fileinfo)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		path = filepath.Dir(path)

		if file, err := match(info.Name()); err == nil {
			key := filepath.Join(path, fmt.Sprintf("%s.%s", file.base, file.ext))

			toRename[key] = append(toRename[key], fileinfo{
				path: path,
				base: file.base,
				ext:  file.ext,
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
		base: pieces[0],
		ext:  ext,
	}, nil
}
