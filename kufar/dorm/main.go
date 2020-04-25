package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	db, err := newDB()
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	c, err := newRedisCacheClient()
	if err != nil {
		log.Fatalf("error connecting to redis: %v", err)
	}

	r := chi.NewRouter()

	r.Middlewares()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/signup", SignUpUser(*db))
	r.Post("/login", LogIn(*db))
	r.Route("/student", func(r chi.Router) {
		r.Post("/", CreateStudent(*db))
		r.Get("/", ReadStudents(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadStudent(*db, c))
			r.Put("/", UpdateStudent(*db, c))
			r.Delete("/", DeleteStudent(*db, c))
		})
	})

	r.Route("/provisions", func(r chi.Router) {
		r.Post("/", CreateProvision(*db))
		r.Get("/", ReadProvisions(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadProvision(*db))
			r.Put("/", UpdateProvision(*db))
			r.Delete("/", DeleteProvision(*db))
		})
	})

	r.Route("/room", func(r chi.Router) {
		r.Post("/", CreateRoom(*db))
		r.Get("/", ReadRooms(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadRoom(*db))
			r.Put("/", UpdateRoom(*db))
			r.Delete("/", DeleteRoom(*db))
		})
	})

	r.Route("/worker", func(r chi.Router) {
		r.Post("/", CreateWorker(*db))
		r.Get("/", ReadWorkers(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadWorker(*db))
			r.Put("/", UpdateWorker(*db))
			r.Delete("/", DeleteWorker(*db))
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
