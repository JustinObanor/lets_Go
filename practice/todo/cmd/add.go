package cmd

import (
	"fmt"
	"strings"

	"github.com/lets_Go/practice/todo/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add ...",
	Short: "Add a new task to your TODO list",
	Long:  "add adds a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {

		value := strings.Join(args, " ")

		if err := db.CreateTask(value); err != nil {
			fmt.Println("couldnt create task", err)
			return
		}

		fmt.Printf("Added \"%s\" to your task list\n", value)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
