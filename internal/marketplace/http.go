package marketplace

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"ufut/internal/showcase"
	funcsUFUT "ufut/lib/funcs"
	structsUFUT "ufut/lib/structs"
)

type Handler struct {
	service  *Service
	showcase *showcase.Handler
}

func NewHandler(srvc *Service) *Handler {
	return &Handler{service: srvc}
}

func (h *Handler) SetShowcase(sc *showcase.Handler) {
	h.showcase = sc
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	handledFuncs := map[string]http.HandlerFunc{
		"POST /api/order/placeOrder":    h.PlaceOrder,
		"POST /api/order/removeOrder":   h.RemoveOrder,
		"GET /api/order/orderStatus":    h.OrderStatus,
		"GET /api/order/userOrders":     h.UserOrders,
		"POST /api/cart/addToCart":      h.AddToCart,
		"POST /api/cart/removeFromCart": h.RemoveFromCart,
		"POST /api/cart/increaseItems":  h.IncreaseItemQuantity,
		"POST /api/cart/decreaseItems":  h.DecreaseItemQuantity,
		"GET /api/cart/listCart":        h.ListCart,
		"POST /api/cart/clearCart":      h.ClearCart,
	}

	for key, val := range handledFuncs {
		mux.HandleFunc(key, funcsUFUT.AuthMiddleware(val))
	}
}

/*
JSON args:

	None

response:

	"status": "ok"
*/
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.PlaceOrder(r.Context(), userID, func(items []string) []bool {
		client := &http.Client{}
		req, err := http.NewRequest("POST", "/api/showcase/reserveItem", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return nil
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		var reqBody struct {
			ItemID     []string `json:"itemID"`
			PassPhrase string   `json:"passphrase"`
		}
		reqBody.ItemID = items
		reqBody.PassPhrase = structsUFUT.PASSPHRASE
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return nil
		}
		req.Body = io.NopCloser(bytes.NewReader(jsonData))
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return nil
		}
		defer resp.Body.Close()
		var respBody struct {
			Successful []bool `json:"successful"`
		}
		json.NewDecoder(resp.Body).Decode(&respBody)
		return respBody.Successful
	}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"orderID": int

response:

	"status": "ok"
*/
func (h *Handler) RemoveOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID int `json:"orderID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.RemoveOrder(r.Context(),
		&structsUFUT.OrderRequestRMP{UserID: userID, OrderID: req.OrderID}, func(items []string) {
			client := &http.Client{}
			req, err := http.NewRequest("POST", "/api/showcase/cancelItemReservation", nil)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			var reqBody struct {
				ItemID     []string `json:"itemID"`
				PassPhrase string   `json:"passphrase"`
			}
			reqBody.ItemID = items
			reqBody.PassPhrase = structsUFUT.PASSPHRASE
			jsonData, err := json.Marshal(reqBody)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			req.Body = io.NopCloser(bytes.NewReader(jsonData))
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
		}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
Query args:

	orderID=int

response:

	"status": any("CREATED", "PREPARING", "DELIVERY", "FINISHED", "CANCELLED")
*/
func (h *Handler) OrderStatus(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	orderID, err := strconv.Atoi(q_vals.Get("orderID"))
	if err != nil {
		http.Error(w, "incorrect orderID", http.StatusBadRequest)
		return
	}
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	req := structsUFUT.OrderRequestRMP{
		OrderID: orderID,
		UserID:  userID,
	}
	if err := h.service.OrderStatus(r.Context(), &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": req.Status})
}

/*
Query args:

	status=string(optional)

resonse:

	{
		"ordersID": [<ints>]
		"statuses": [<strings>]
	}
*/
func (h *Handler) UserOrders(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	status := q_vals.Get("status")
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	resp, err := h.service.UserOrders(r.Context(), &structsUFUT.OrderRequestRMP{UserID: userID, Status: status})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"itemID": string
	"quantity":int

response:

	"status": "ok"
*/
func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID   string `json:"itemID"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.AddToCart(r.Context(), &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID, Quantity: req.Quantity}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"itemID":string

response:

	"status": "ok"
*/
func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID string `json:"itemID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.RemoveFromCart(r.Context(), &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"itemID": string
	"quantity":int

response:

	"status": "ok"
*/
func (h *Handler) IncreaseItemQuantity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID   string `json:"itemID"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.IncreaseItemQuantity(r.Context(), &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID, Quantity: req.Quantity}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"itemID": string
	"quantity":int

response:

	"status": "ok"
*/
func (h *Handler) DecreaseItemQuantity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID   string `json:"itemID"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.DecreaseItemQuantity(r.Context(), &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID, Quantity: req.Quantity}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
Query args:

	None

response
*/
func (h *Handler) ListCart(w http.ResponseWriter, r *http.Request) {
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	resp, err := h.service.ListCart(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var t_strct struct {
		T1 []string `json:"itemsID"`
		T2 []int    `json:"quantities"`
	}
	t_strct.T1 = resp.ItemsID
	t_strct.T2 = resp.Quantities
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(t_strct)
}

/*
JSON args:

	None

response:

	"status": "ok"
*/
func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.ClearCart(r.Context(), userID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
