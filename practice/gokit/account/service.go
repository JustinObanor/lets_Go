package account

import "context"

//methods to be exposed to transport & use the interface to implement the bussiness logic
//theyll be exposed from our microservice

type Service interface {
	CreateUser(ctx context.Context, email string, password string) (string, error)
	GetUser(ctx context.Context, id string) (string, error)
}
