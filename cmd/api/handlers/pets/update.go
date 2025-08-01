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
	Name      string                `json:"name"`
	Category  responses.PetCategory `json:"category"`
	PhotoUrls []string              `json:"photo_urls"`
	Tags      []responses.PetTag    `json:"tags"`
	Status    string                `json:"status"`
}
type UpdateResponse struct {
	Body responses.Envelope[responses.Pet]
}

func (h *Handlers) Update(_ context.Context, _ *UpdateRequest) (*UpdateResponse, error) {
	return &UpdateResponse{}, nil
}
