package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var bucket = "tasks"
var db *bolt.DB

type task struct {
	ID    string
	Value string
}

func NewDB() error {
	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func CreateTask(value string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		if b.Stats().KeyN == 0 {
			if err := b.SetSequence(0); err != nil {
				return err
			}
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		return b.Put([]byte(fmt.Sprint(id)), []byte(value))
	})
}

func RemoveTask(id string) (string, error) {
	var task []byte
	return string(task), db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		if task = b.Get([]byte(id)); task != nil {
			return b.Delete([]byte(id))
		}

		return nil
	})
}

func GetTasks() ([]task, error) {
	var tasks []task

	return tasks, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, task{
				ID:    string(k),
				Value: string(v),
			})
		}
		return nil
	})
}
