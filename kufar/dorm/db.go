package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

var host = getenv("PSQL_HOST", "db")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWDcas", "postgres")
var dbname = getenv("PSQL_DB_NAME", "dorm")

//Database ...
type Database struct {
	db *sql.DB
}

type rediscache struct {
	redis *redis.Client
}

// newRedisCacheClient ...
func newRedisCacheClient() (*rediscache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "db_redis:6379",
		Password: os.Getenv("Password"),
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Println("redis connected")

	return &rediscache{
		redis: client,
	}, nil
}

func newDB() (*Database, error) {
	var err error

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s
	password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("postgres connected")

	return &Database{
		db: db,
	}, nil
}

//Close closes the conn
func (d Database) Close() error {
	return d.db.Close()
}
