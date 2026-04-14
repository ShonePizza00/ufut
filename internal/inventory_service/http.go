package inventory_service

import "net/http"

type Handler struct {
	service *Service
}

func NewHandler(srvc *Service) *Handler {
	return &Handler{service: srvc}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	handledFuncs := map[string]http.HandlerFunc{
		// "POST /api/items/reserve":  h.ReserveItem,
		// "POST /api/items/cancel":   h.CancelItemReservation,
		"POST /api/items/increase": h.IncreaseQuantity,
		"POST /api/items/decrease": h.DecreaseQuantity,
		"GET /api/items/quantity":  h.ItemQuantity,
	}

	for key, val := range handledFuncs {
		mux.HandleFunc(key, val)
	}
}

// func (h *Handler) ReserveItem(w http.ResponseWriter, r *http.Request) {

// }

// func (h *Handler) CancelItemReservation(w http.ResponseWriter, r *http.Request) {

// }

func (h *Handler) IncreaseQuantity(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DecreaseQuantity(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ItemQuantity(w http.ResponseWriter, r *http.Request) {

}
