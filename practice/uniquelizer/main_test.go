package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var test = `1
2
3
3
4
5
5
6`

var res = `1
2
3
4
5
6
`

func TestOK(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(test))
	out := new(bytes.Buffer)

	err := uniquerlizer(in, out)
	if err != nil {
		t.Error("tests failed (error found)")
	}

	buf := out.String()
	if buf != res {
		t.Errorf("\ntests failed (results dont match)\n got %v expected %v\n", buf, res)
	}
}

var unSorted = `1
2
1`

func TestForUnsorted(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(unSorted))
	out := new(bytes.Buffer)

	err := uniquerlizer(in, out)
	if err == nil {
		t.Errorf("tests failed: %v", err)
	}

}
