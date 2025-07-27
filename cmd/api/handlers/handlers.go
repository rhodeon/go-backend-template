package handlers

import "github.com/rhodeon/go-backend-template/cmd/api/internal"

type Handlers struct {
	*baseHandler
	Users *usersHandler
}

func New(app *internal.Application) *Handlers {
	return &Handlers{
		newBaseHandler(app),
		newUsersHandler(app),
	}
}
