package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := uniquerlizer(os.Stdin, os.Stdout); err != nil {
		panic(err.Error())
	}
}

func uniquerlizer(r io.Reader, w io.Writer) error {
	b := bufio.NewScanner(r)

	var prev string
	for b.Scan() {
		txt := b.Text()

		if prev == txt {
			continue
		}

		if prev > txt {
			return errors.New("file not sorted")
		}

		prev = txt

		fmt.Fprintln(w, txt)
	}
	return nil
}
