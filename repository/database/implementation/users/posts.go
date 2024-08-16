package users

import "github.com/rhodeon/go-backend-template/repository/database/common_models"

func (p Post) Commonize() common_models.Post {
	return common_models.Post(p)
}
