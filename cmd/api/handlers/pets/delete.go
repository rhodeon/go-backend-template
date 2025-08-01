package pets

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type DeleteRequest struct {
	PetId int `path:"pet_id"`
}

type DeleteResponse struct {
	Body responses.Envelope[responses.SuccessMessage]
}

func (h *Handlers) delete(_ context.Context, _ *DeleteRequest) (*DeleteResponse, error) {
	return &DeleteResponse{
		Body: responses.Success[responses.SuccessMessage]("Pet deleted"),
	}, nil
}
