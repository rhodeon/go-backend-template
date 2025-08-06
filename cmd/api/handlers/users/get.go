package users

import (
	"context"

	apierrors "github.com/rhodeon/go-backend-template/cmd/api/errors"
	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"
	domainerrors "github.com/rhodeon/go-backend-template/domain/errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/pkg/errors"
)

type GetRequest struct {
	UserId int64 `path:"user_id"`
}

type GetResponse struct {
	Body responses.Envelope[responses.User]
}

func (h *Handlers) get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	dbTx, commit, rollback, err := h.app.Db.BeginTx(ctx)
	if err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}
	defer rollback(ctx)

	user, err := h.app.Services.User.GetById(ctx, dbTx, req.UserId)
	if err != nil {
		var errRecordNotFound *domainerrors.RecordNotFoundError

		switch {
		case errors.As(err, &errRecordNotFound):
			if errRecordNotFound.Entity == "user" {
				return nil, huma.Error404NotFound("user not found")
			}

			fallthrough

		default:
			return nil, apierrors.UntypedError(ctx, err)
		}
	}

	respData := responses.NewUser.FromDomainUser(user)

	if err := commit(ctx); err != nil {
		return nil, apierrors.UntypedError(ctx, err)
	}

	return &GetResponse{
		Body: responses.Success(respData),
	}, nil
}
