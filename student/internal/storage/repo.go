package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/lets_Go/student/internal/student"
)

var RepoErr = errors.New("Unable to handle Repo Request")

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) student.Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "postgres"),
	}
}

func (repo *repo) CreateStudent(ctx context.Context, student student.Student) error {
	sql := `insert into students (id, firstname, lastname, room) values($1, $2, $3, $4)`

	if student.FirstName == "" || student.LastName == "" {
		return RepoErr
	}

	if _, err := repo.db.ExecContext(ctx, sql, student.ID, student.FirstName, student.LastName, student.StudentRoom); err != nil {
		return err
	}
	return nil
}

func (repo *repo) GetStudent(ctx context.Context, id string) (student.Student, error) {
	var student student.Student

	sql := `select * from students where id = $1`

	if err := repo.db.QueryRow(sql, id).Scan(&student.ID, &student.FirstName, &student.LastName, &student.StudentRoom); err != nil {
		return student, err
	}

	return student, nil
}

func (repo *repo) UpdateStudent(ctx context.Context, id string, student student.Student) error {
	sql := `update students set firstname = $1, lastname = $2, room = $3 where id = $4`

	if student.FirstName == "" || student.LastName == "" {
		return RepoErr
	}

	if _, err := repo.db.ExecContext(ctx, sql, student.FirstName, student.LastName, student.StudentRoom, student.ID); err != nil {
		return err
	}

	return nil
}

func (repo *repo) DeleteStudent(ctx context.Context, id string) error {
	sql := `delete from students where id = $1`

	if _, err := repo.db.ExecContext(ctx, sql, id); err != nil {
		return err
	}

	return nil
}
