package pets

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
)

type ListRequest struct {
	Status string   `query:"status"`
	Tags   []string `query:"tags"`
}

type ListResponse struct {
	Body responses.EnvelopeWithMetadata[[]responses.Pet, *responses.Pagination]
}

func (h *Handlers) list(_ context.Context, _ *ListRequest) (*ListResponse, error) {
	return &ListResponse{
		Body: responses.SuccessWithMetadata(
			[]responses.Pet{},
			responses.CalculatePagination(0, 0, 0),
		),
	}, nil
}
