package transport

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/lets_Go/student/internal/student"
)

func NewCreateStudentEndpoint(svc student.Service) http.HandlerFunc {
	return httptransport.NewServer(
		makeCreateStudentEndpoint(svc),
		decodeCreateUpdateStudentRequest,
		encodeResponse,
	).ServeHTTP
}

func NewGetStudentEndpoint(svc student.Service) http.HandlerFunc {
	return httptransport.NewServer(
		makeGetStudentEndpoint(svc),
		decodeGetDeleteStudentRequest,
		encodeResponse,
	).ServeHTTP
}

func NewUpdateStudentEndpoint(svc student.Service) http.HandlerFunc {
	return httptransport.NewServer(
		makeUpdateStudentEndpoint(svc),
		decodeCreateUpdateStudentRequest,
		encodeResponse,
	).ServeHTTP
}

func NewDeleteStudentEndpoint(svc student.Service) http.HandlerFunc {
	return httptransport.NewServer(
		makeDeleteStudentEndpoint(svc),
		decodeGetDeleteStudentRequest,
		encodeResponse,
	).ServeHTTP
}
