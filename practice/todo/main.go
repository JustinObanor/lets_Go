package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const (
	bucket = "TODO"
)

type database struct {
	db *bolt.DB
}

func newDatabse() *database {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &database{
		db: db,
	}
}

func main() {
	db := newDatabse()
	defer db.db.Close()

	db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})

	var rootCmd = &cobra.Command{
		Use:   "task [command]",
		Short: "task is a CLI for managing your TODOs.",
	}

	rootCmd.AddCommand(db.addCmd(), db.doCmd(), db.listCmd())
	rootCmd.Execute()
}

func (db *database) addCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "add ...",
		Short: "Add a new task to your TODO list",
		Long:  "add adds a new task to your TODO list",
		Run: func(cmd *cobra.Command, args []string) {

			db.db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bucket))

				if b.Stats().KeyN == 0 {
					if err := b.SetSequence(0); err != nil {
						log.Println(err)
					}
				}

				id, err := b.NextSequence()
				if err != nil {
					log.Println(err)
				}

				return b.Put([]byte(fmt.Sprint(id)),
					[]byte(strings.Join(args, " ")))
			})

			fmt.Printf("Added \"%s\" to your task list\n", strings.Join(args, " "))
		},
	}
}

func (db *database) doCmd() *cobra.Command {

	return &cobra.Command{
		Use:   "do [taskID]",
		Short: "Mark a task on your TODO list as complete",
		Long:  "do marks a task on your TODO list as complete and removes from list",
		Run: func(cmd *cobra.Command, args []string) {

			db.db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bucket))

				if task := b.Get([]byte(args[0])); task != nil {
					fmt.Printf("You have completed the \"%s\" task, id %s.\n", string(task), args[0])

					b.Delete([]byte(args[0]))

				} else {
					fmt.Println("No such task")
				}
				return nil
			})
		},
	}
}

func (db *database) listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [taskID] | [all]",
		Short: "List all of your incomplete tasks",
		Long:  "list lists all of your incomplete tasks in your TODO list",
		Run: func(cmd *cobra.Command, args []string) {

			db.db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bucket))

				c := b.Cursor()

				fmt.Println("You have the following tasks:")
				for k, v := c.First(); k != nil; k, v = c.Next() {
					fmt.Printf("%s. %s.\n", k, v)
				}
				return nil
			})
		},
	}
}
