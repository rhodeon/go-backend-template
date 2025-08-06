package responses

import (
	"time"

	"github.com/rhodeon/go-backend-template/domain/models"
)

var NewUser = new(User)

type User struct {
	Id        int64     `json:"id" required:"true" example:"1"`
	Username  string    `json:"username" required:"true" example:"johndoe"`
	FirstName string    `json:"first_name" required:"true" example:"John"`
	LastName  string    `json:"last_name" required:"true" example:"Doe"`
	Email     string    `json:"email" required:"true" example:"johndoe@example.com"`
	Phone     string    `json:"phone" required:"false"`
	CreatedAt time.Time `json:"created_at" required:"true"`
}

func (u *User) FromDomainUser(domainUser models.User) User {
	return User{
		Id:        domainUser.Id,
		Username:  domainUser.Username,
		FirstName: domainUser.FirstName,
		LastName:  domainUser.LastName,
		Email:     domainUser.Email,
		Phone:     domainUser.PhoneNumber,
		CreatedAt: domainUser.CreatedAt,
	}
}
