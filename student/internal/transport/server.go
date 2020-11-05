package transport

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(ctx context.Context, endpoint Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/student").Handler(httptransport.NewServer(
		endpoint.CreateStudent,
		decodeCreateUpdateStudentRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/student/{id}").Handler(httptransport.NewServer(
		endpoint.GetStudent,
		decodeGetDeleteStudentRequest,
		encodeResponse,
	))

	r.Methods("PUT").Path("/student/{id}").Handler(httptransport.NewServer(
		endpoint.UpdateStudent,
		decodeCreateUpdateStudentRequest,
		encodeResponse,
	))

	r.Methods("DELETE").Path("/student/{id}").Handler(httptransport.NewServer(
		endpoint.DeleteStudent,
		decodeGetDeleteStudentRequest,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
