package cmd

import (
	"fmt"

	"github.com/lets_Go/practice/todo/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [taskID] | [all]",
	Short: "List all of your incomplete tasks",
	Long:  "list lists all of your incomplete tasks in your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.GetTasks()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("You have the following tasks:")
		for _, task := range tasks {
			fmt.Printf("%s. %s\n", task.ID, task.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
