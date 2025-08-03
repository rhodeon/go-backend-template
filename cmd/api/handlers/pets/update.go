package pets

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type UpdateRequest struct {
	Body  UpdateRequestBody
	PetId int `path:"pet_id"`
}

type UpdateRequestBody struct {
	Name      string                `json:"name" required:"false"`
	Category  responses.PetCategory `json:"category" required:"false"`
	PhotoUrls []string              `json:"photo_urls" required:"false"`
	Tags      []responses.PetTag    `json:"tags" required:"false"`
	Status    string                `json:"status" required:"false"`
}
type UpdateResponse struct {
	Body responses.Envelope[responses.Pet]
}

func (h *Handlers) update(_ context.Context, _ *UpdateRequest) (*UpdateResponse, error) {
	return &UpdateResponse{}, nil
}
