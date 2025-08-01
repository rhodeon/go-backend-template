package pets

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type CreateRequest struct {
	Body CreateRequestBody
}

type CreateRequestBody struct {
	Name      string                `json:"name"`
	Category  responses.PetCategory `json:"category"`
	PhotoUrls []string              `json:"photo_urls"`
	Tags      []responses.PetTag    `json:"tags"`
	Status    string                `json:"status"`
}

type CreateResponse struct {
	Body responses.Envelope[responses.Pet]
}

func (h *Handlers) Create(_ context.Context, _ *CreateRequest) (*CreateResponse, error) {
	return &CreateResponse{}, nil
}
