package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lets_Go/student/internal/storage"
	"github.com/lets_Go/student/internal/student"
	"github.com/lets_Go/student/internal/transport"
)

var host = getenv("PSQL_HOST", "studentdb")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWDcas", "1999")
var dbname = getenv("PSQL_DB_NAME", "student")

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error

		db, err = sql.Open("postgres", fmt.Sprintf(`host=%s port=%s user=%s
		password=%s dbname=%s sslmode=disable`,
			host, port, user, password, dbname))
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	flag.Parse()

	var svc student.Service
	{
		repository := storage.NewRepo(db, logger)

		svc = student.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		r.Post("/student", transport.NewCreateStudentEndpoint(svc))
		r.Get("/student/{id}", transport.NewGetStudentEndpoint(svc))
		r.Put("/student/{id}", transport.NewUpdateStudentEndpoint(svc))
		r.Delete("/student/{id}", transport.NewDeleteStudentEndpoint(svc))
	})

	go func() {
		fmt.Println("listening on port", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, r)
	}()

	level.Error(logger).Log("exit", <-errs)
}
