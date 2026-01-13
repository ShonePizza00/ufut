package httpserver

import (
	"context"
	"net/http"
	"ufut/internal/auth"
	"ufut/internal/marketplace"
	"ufut/internal/showcase"
)

type Services struct {
	Service_Auth *auth.Service
	Service_MP   *marketplace.Service
	Service_SC   *showcase.Service
}

func AddRoutes(ctx context.Context, mux *http.ServeMux, services *Services) {
	handler_Auth := auth.NewHandler(services.Service_Auth)
	handler_MP := marketplace.NewHandler(services.Service_MP)
	handler_SC := showcase.NewHandler(services.Service_SC)

	handler_MP.SetShowcase(handler_SC)

	auth.RegisterRoutes(mux, handler_Auth)
	marketplace.RegisterRoutes(mux, handler_MP)
	showcase.RegisterRoutes(mux, handler_SC)
}
