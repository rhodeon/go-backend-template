package responses

import "time"

type User struct {
	Id        int64     `json:"id" required:"true" example:"1"`
	Username  string    `json:"username" required:"true" example:"johndoe"`
	FirstName string    `json:"first_name" required:"true" example:"John"`
	LastName  string    `json:"last_name" required:"true" example:"Doe"`
	Email     string    `json:"email" required:"true" example:"johndoe@example.com"`
	Phone     string    `json:"phone" required:"false"`
	CreatedAt time.Time `json:"created_at" required:"true"`
}
