package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/MyWorkSpace/lets_Go/practice/protocols/todo"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()

	//NArg is the number of arguments remaining after flags have been processed.
	//helps us to only run with flags
	//The subcommand is expected as the first argument to the program.
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missiing subcommand : list or add")
		//Exit causes the current program to exit with the given status code
		//0 == success, !0 == error
		os.Exit(1)
	}
	var err error
	//Arg(0) is the first remaining argument after flags have been processed
	//Check which subcommand is invoked.
	//For every subcommand, we parse its own flags and have access to trailing positional arguments.
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list()
	case "add":
		//Args returns the non-flag command-line arguments.
		err = add(strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type length int64

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)

var endianness = binary.LittleEndian

func add(text string) error {
	task := &todo.Task{
		Text: text,
		Done: false,
	}
	//marshalling it and turning it into the format
	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("coulnt not encode task: %v", err)
	}
	//if we ere able to encode it, the we write to the file
	//write only, create if it didnt exist, if it exists, we append
	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("coulnt not encode task %s : %v", dbPath, err)
	}

	/*
		//write the length b4 we encode the messagge
		//encode the len of 1nt64 of msg
		if err := gob.NewEncoder(f).Encode(int64(len(b))); err != nil {
			return fmt.Errorf("could not encode length of message: %v", err)
		}
	*/

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}

	//otherwise, we write to the file
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	//check error since we writing
	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file %s : %v", dbPath, err)
	}
	return nil
}

func list() error {
	//read file
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read %s : %v", dbPath, err)
	}

	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < sizeOfLength {
			return fmt.Errorf("remaining odd %d bytes, what to do?", len(b))
		}

		var l length
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]
		//get first 8 bytes

		var task todo.Task
		//unmarshall so we can read it
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return fmt.Errorf("coulnd not read task: %v", err)
		}
		b = b[l:]

		//if true
		if task.Done {
			fmt.Printf("[👍]")
		} else {
			fmt.Printf("😱")
		}
		//print out data we unmarshelled
		fmt.Printf("%s\n", task.Text)
	}
}
