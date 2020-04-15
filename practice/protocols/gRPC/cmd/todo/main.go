package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MyWorkSpace/lets_Go/practice/protocols/gRPC/todo"
	grpc "google.golang.org/grpc"
)

//client thayt connect through gRPC

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

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to backend: %v", err)
	}
	client := todo.NewTasksClient(conn)

	//Arg(0) is the first remaining argument after flags have been processed
	//Check which subcommand is invoked.
	//For every subcommand, we parse its own flags and have access to trailing positional arguments.
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "add":
		//Args returns the non-flag command-line arguments.
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(ctx context.Context, client todo.TasksClient, text string) error {
	_, err := client.Add(ctx, &todo.Text{Text: text})
	if err != nil {
		return fmt.Errorf("couldnt not add task through the backend: %v", err)
	}
	fmt.Println("task added successfully")

	return nil
}

func list(ctx context.Context, client todo.TasksClient) error {
	l, err := client.List(ctx, &todo.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch tassks: %v", err)
	}
	for _, t := range l.Tasks {
		if t.Done {
			fmt.Printf("👍")
		} else {
			fmt.Printf("😱")
		}
		fmt.Printf("%s\n", t.Text)
	}
	return nil
}
