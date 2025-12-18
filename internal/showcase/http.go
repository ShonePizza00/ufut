package showcase

import "context"

type Handler struct {
	service *Service
	ctx     context.Context
}

func NewHandler(ctx context.Context, srvc *Service) *Handler {
	return &Handler{service: srvc, ctx: ctx}
}
