package marketplace

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	structsUFUT "ufut/lib"
)

var (
	ErrIncorrectToken = errors.New("incorrect token")
)

type Handler struct {
	service *Service
	ctx     context.Context
}

func NewHandler(ctx context.Context, srvc *Service) *Handler {
	return &Handler{service: srvc, ctx: ctx}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /api/placeorder", h.PlaceOrder)
	mux.HandleFunc("POST /api/removeorder", h.RemoveOrder)
	mux.HandleFunc("GET /api/OrderStatus", h.OrderStatus)
	mux.HandleFunc("GET /api/userorders", h.UserOrders)
	mux.HandleFunc("POST /api/addtocart", h.AddToCart)
	mux.HandleFunc("POST /api/removefromcart", h.RemoveFromCart)
	mux.HandleFunc("POST /api/increaseitems", h.IncreaseItemQuantity)
	mux.HandleFunc("POST /api/decreaseitems", h.DecreaseItemQuantity)
	mux.HandleFunc("GET /api/listcart", h.ListCart)
	mux.HandleFunc("POST /api/clearcart", h.ClearCart)
}

func useridByToken(token string) (string, error) {
	req, err := http.NewRequest("GET", "/api/verifyTokenUser", nil)
	if err != nil {
		return "", err
	}
	query := req.URL.Query()
	query.Set("token", token)
	query.Set("passphrase", structsUFUT.PASSPHRASE)
	req.URL.RawQuery = query.Encode()
	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var id_resp struct {
		UserID string `json:"userID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&id_resp); err != nil {
		return "", err
	}
	return id_resp.UserID, nil
}

/*
JSON args:

	"token": string

response:

	"status": "ok"
*/
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var token structsUFUT.TokenResponse
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(token.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.PlaceOrder(h.ctx, userID); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"token": string
	"orderID": int

response:

	"status": "ok"
*/
func (h *Handler) RemoveOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token   string `json:"token"`
		OrderID int    `json:"orderID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(req.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.RemoveOrder(h.ctx,
		&structsUFUT.OrderRequestRMP{UserID: userID, OrderID: req.OrderID}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
Query args:

	token=string
	orderID=int

response:

	"status": any("CREATED", "PREPARING", "DELIVERY", "FINISHED", "CANCELLED")
*/
func (h *Handler) OrderStatus(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	token := q_vals.Get("token")
	orderID, err := strconv.Atoi(q_vals.Get("orderID"))
	if err != nil {
		http.Error(w, "incorrect orderID", http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	req := structsUFUT.OrderRequestRMP{
		OrderID: orderID,
		UserID:  userID,
	}
	if err := h.service.OrderStatus(h.ctx, &req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": req.Status})
}

/*
Query args:

	token=string
	status=string(optional)

resonse:

	{
		"ordersID": [<ints>]
		"statuses": [<strings>]
	}
*/
func (h *Handler) UserOrders(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	token := q_vals.Get("token")
	status := q_vals.Get("status")
	userID, err := useridByToken(token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	resp, err := h.service.UserOrders(h.ctx, &structsUFUT.OrderRequestRMP{UserID: userID, Status: status})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"token":string
	"itemID": string
	"quantity":int

response:

	"status": "ok"
*/
func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		ItemID   string `json:"itemID"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(req.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.AddToCart(h.ctx, &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID, Quantity: req.Quantity}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"token":string
	"itemID":string

response:

	"status": "ok"
*/
func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token  string `json:"token"`
		ItemID string `json:"itemID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(req.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.RemoveFromCart(h.ctx, &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"token":string
	"itemID": string
	"quantity":int

response:

	"status": "ok"
*/
func (h *Handler) IncreaseItemQuantity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		ItemID   string `json:"itemID"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(req.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.IncreaseItemQuantity(h.ctx, &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID, Quantity: req.Quantity}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"token":string
	"itemID": string
	"quantity":int

response:

	"status": "ok"
*/
func (h *Handler) DecreaseItemQuantity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		ItemID   string `json:"itemID"`
		Quantity int    `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(req.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.DecreaseItemQuantity(h.ctx, &structsUFUT.ItemRequestRMP{
		UserID: userID, ItemID: req.ItemID, Quantity: req.Quantity}); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
Query args:

	token=string

response
*/
func (h *Handler) ListCart(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	token := q_vals.Get("token")
	userID, err := useridByToken(token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	resp, err := h.service.ListCart(h.ctx, userID)
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

	"token": string

response:

	"status": "ok"
*/
func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	var token structsUFUT.TokenResponse
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	userID, err := useridByToken(token.Token)
	if err != nil {
		http.Error(w, ErrIncorrectToken.Error(), http.StatusForbidden)
		return
	}
	if err := h.service.ClearCart(h.ctx, userID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
