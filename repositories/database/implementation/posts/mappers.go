package posts

import "github.com/rhodeon/go-backend-template/repositories/database/commonmodels"

func (p Post) Commonize() commonmodels.Post {
	return commonmodels.Post(p)
}
