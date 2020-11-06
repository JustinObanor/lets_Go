package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/lets_Go/student/internal/student"
)

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
