package student

import "context"

type Student struct {
	ID          string `json:"id,omitempty"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	StudentRoom int    `json:"room"`
}

type Repository interface {
	CreateStudent(ctx context.Context, student Student) error
	GetStudent(ctx context.Context, id string) (Student, error)
	UpdateStudent(ctx context.Context, id string, student Student) error
	DeleteStudent(ctx context.Context, id string) error
}
