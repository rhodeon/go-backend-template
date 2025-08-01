package pets

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type GetRequest struct {
	PetId int `path:"pet_id"`
}

type GetResponse struct {
	Body responses.Envelope[responses.Pet]
}

func (h *Handlers) get(_ context.Context, _ *GetRequest) (*GetResponse, error) {
	return &GetResponse{
		Body: responses.Success(responses.Pet{}),
	}, nil
}
