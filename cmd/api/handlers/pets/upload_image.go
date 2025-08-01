package pets

import (
	"context"

	"github.com/rhodeon/go-backend-template/cmd/api/models/responses"

	"github.com/danielgtaylor/huma/v2"
)

type UploadImageRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		AvatarFile huma.FormFile `form:"image" required:"true" doc:"max size: 1 MB; allowed file format: png, jpg, jpeg"`
	}]

	PetId int `path:"pet_id"`
}

type (
	UploadImageRequestBody struct{}
	UploadImageResponse    struct {
		Body responses.Envelope[responses.SuccessMessage]
	}
)

func (h *Handlers) UploadImage(_ context.Context, _ *UpdateRequest) (*UpdateResponse, error) {
	return &UpdateResponse{}, nil
}
