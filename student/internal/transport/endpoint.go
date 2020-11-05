package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/lets_Go/student/internal/student"
)

type Endpoints struct {
	CreateStudent endpoint.Endpoint
	GetStudent    endpoint.Endpoint
	UpdateStudent endpoint.Endpoint
	DeleteStudent endpoint.Endpoint
}

func MakeEndpoints(s student.Service) Endpoints {
	return Endpoints{
		CreateStudent: makeCreateStudentEndpoint(s),
		GetStudent:    makeGetStudentEndpoint(s),
		UpdateStudent: makeUpdateStudentEndpoint(s),
		DeleteStudent: makeDeleteStudentEndpoint(s),
	}
}

func makeCreateStudentEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUpdateStudentRequest)
		ok, err := s.CreateStudent(ctx, req.FirstName, req.LastName, req.StudentRoom)

		return CreateUpdateStudentResponse{Ok: ok}, err
	}
}

func makeGetStudentEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetStudentRequest)
		stud, err := s.GetStudent(ctx, req.ID)

		return GetStudentResponse{
			ID:          stud.ID,
			FirstName:   stud.FirstName,
			LastName:    stud.LastName,
			StudentRoom: stud.StudentRoom,
		}, err
	}
}

func makeUpdateStudentEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateUpdateStudentRequest)
		ok, err := s.UpdateStudent(ctx, req.ID, req.FirstName, req.LastName, req.StudentRoom)

		return CreateUpdateStudentResponse{Ok: ok}, err
	}
}

func makeDeleteStudentEndpoint(s student.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetStudentRequest)
		ok, err := s.DeleteStudent(ctx, req.ID)

		return DeleteStudentResponse{Ok: ok}, err
	}
}
