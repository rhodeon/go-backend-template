package orders

import "context"

type GetByIdRequest struct {
	Id string `path:"id"`
}

func (h *Handlers) GetById(ctx context.Context, req *GetByIdRequest) (*struct{}, error) {
	return nil, nil
}
