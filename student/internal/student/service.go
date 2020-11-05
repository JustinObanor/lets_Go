package student

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

type Service interface {
	CreateStudent(ctx context.Context, firstName, lastName string, studRoom int) (string, error)
	GetStudent(ctx context.Context, id string) (Student, error)
	UpdateStudent(ctx context.Context, id, firstName, lastName string, studRoom int) (string, error)
	DeleteStudent(ctx context.Context, id string) (string, error)
}

type service struct {
	repostory Repository
	logger    log.Logger
}

func NewService(rep Repository, logger log.Logger) Service {
	return &service{
		repostory: rep,
		logger:    logger,
	}
}

func (s service) CreateStudent(ctx context.Context, firstName, lastName string, studRoom int) (string, error) {
	logger := log.With(s.logger, "method", "CreateStudent")

	uuid, _ := uuid.NewV4()
	id := uuid.String()

	student := Student{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		StudentRoom: studRoom,
	}

	if err := s.repostory.CreateStudent(ctx, student); err != nil {
		level.Error(logger).Log("err", err)
		return "", nil
	}

	logger.Log("create student", id)

	return "success", nil
}

func (s service) GetStudent(ctx context.Context, id string) (Student, error) {
	logger := log.With(s.logger, "method", "GetStudent")

	student, err := s.repostory.GetStudent(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return Student{}, err
	}

	logger.Log("get student", id)

	return student, nil
}

func (s service) UpdateStudent(ctx context.Context, id string, firstName, lastName string, studRoom int) (string, error) {
	logger := log.With(s.logger, "method", "UpdateStudent")

	student := Student{
		ID: id,
		FirstName:   firstName,
		LastName:    lastName,
		StudentRoom: studRoom,
	}

	if err := s.repostory.UpdateStudent(ctx, id, student); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("update student", id)

	return "success", nil
}

func (s service) DeleteStudent(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteStudent")

	if err := s.repostory.DeleteStudent(ctx, id); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}

	logger.Log("delete student", id)

	return "success", nil
}
