package cmd

import (
	"fmt"

	"github.com/lets_Go/practice/todo/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do [taskID]",
	Short: "Mark a task on your TODO list as complete",
	Long:  "do marks a task on your TODO list as complete and removes from list",
	Run: func(cmd *cobra.Command, args []string) {
		task, err := db.RemoveTask(args[0])
		if err != nil {
			fmt.Println("no such task", err)
			return
		}

		fmt.Printf("You have completed the \"%s\" task.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
