package pets

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	"github.com/rhodeon/go-backend-template/domain"
)

type CreateRequest struct {
	Body CreateRequestBody
}

type CreateRequestBody struct {
	Name       string   `json:"name"`
	CategoryId int64    `json:"category_id"`
	PhotoUrls  []string `json:"photo_urls"`
	Tags       []string `json:"tags"`
	Status     string   `json:"status"`
}

type CreateResponse struct {
	Body responses.Envelope[responses.Pet]
}

func (h *Handlers) create(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	dbTx, commit, rollback, err := h.app.Db.BeginTx(ctx)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	createdPet, err := h.app.Services.Pet.Create(ctx, dbTx, domain.Pet{
		Name: req.Body.Name,
		Category: domain.PetCategory{
			Id: req.Body.CategoryId,
		},
		Status:    req.Body.Status,
		ImageUrls: req.Body.PhotoUrls,
	})
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	respData := responses.NewPet.FromDomainPet(createdPet)

	return &CreateResponse{
		Body: responses.Success(respData),
	}, nil
}
