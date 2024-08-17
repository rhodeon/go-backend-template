package requests

type CreateUserRequest struct {
	Body CreateUserRequestBody
}

type CreateUserRequestBody struct {
	Username string `json:"username" required:"true"`
	Email    string `json:"email" required:"true" format:"email"`
}

type GetUserRequest struct {
	Id int `json:"id" path:"id"`
}
