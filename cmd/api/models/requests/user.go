package requests

type UsersCreateRequest struct {
	Body UsersCreateRequestBody
}

type UsersCreateRequestBody struct {
	Username string `json:"username" required:"true" minLength:"1"`
	Email    string `json:"email" required:"true" format:"email" minLength:"1"`
}

type UsersGetByIdRequest struct {
	Id int `json:"id" path:"id"`
}
