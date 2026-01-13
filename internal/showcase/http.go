package showcase

import (
	"encoding/json"
	"net/http"
	"strconv"
	funcsUFUT "ufut/lib/funcs"
	structsUFUT "ufut/lib/structs"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(srvc *Service) *Handler {
	return &Handler{service: srvc}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	handledFuncs := map[string]http.HandlerFunc{
		"GET /api/user/categories":    h.Categories,
		"GET /api/user/itemsByParams": h.ItemsByParams,
		"GET /api/user/itemByItemID":  h.ItemByItemID,
		"POST /api/staff/createItem":  h.CreateItem,
		"POST /api/staff/deleteItem":  h.DeleteItem,

		"POST /api/showcase/reserveItem":           h.ReserveItem,
		"POST /api/showcase/cancelItemReservation": h.CancelItemReservation,
	}

	for key, val := range handledFuncs {
		mux.HandleFunc(key, funcsUFUT.AuthMiddleware(val))
	}

	// mux.HandleFunc("GET /api/user/categories", h.Categories)
	// mux.HandleFunc("GET /api/user/itemsByParams", h.ItemsByParams)
	// mux.HandleFunc("GET /api/user/itemByItemID", h.ItemByItemID)
	// mux.HandleFunc("POST /api/staff/createItem", h.CreateItem)
	// mux.HandleFunc("POST /api/staff/deleteItem", h.DeleteItem)
}

// func staffidByToken(token string) (string, error) {
// 	req, err := http.NewRequest("GET", "/api/verifyTokenStaff", nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	query := req.URL.Query()
// 	query.Set("token", token)
// 	query.Set("passphrase", structsUFUT.PASSPHRASE)
// 	req.URL.RawQuery = query.Encode()
// 	clt := http.Client{}
// 	resp, err := clt.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	var id_resp struct {
// 		StaffID string `json:"staffID"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&id_resp); err != nil {
// 		return "", err
// 	}
// 	return id_resp.StaffID, nil
// }

/*
Query args:

	None

Response:

	"categories": []string
*/
func (h *Handler) Categories(w http.ResponseWriter, r *http.Request) {
	var resp struct {
		Categories []string `json:"categories"`
	}
	res, err := h.service.Categories(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	resp.Categories = res
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
Query args:

	category: always (specifies which category to search in)
	price: TODO
	startindex: optional, 0 if not provided (offset from begging)
	count: optional, 10 if not provided (number of items in response)
	orderby: "asc" or "desc". optional, "desc" if not provided. (specifies order)

resp:

	itemsID: array of <string>ItemID
*/
func (h *Handler) ItemsByParams(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	var params structsUFUT.ItemsRequestRSC
	params.Category = q_vals.Get("category")
	params.Price, _ = strconv.Atoi(q_vals.Get("price"))
	params.StartIndex, _ = strconv.Atoi(q_vals.Get("startindex"))
	params.Count, _ = strconv.Atoi(q_vals.Get("count"))
	params.OrderBy = q_vals.Get("orderby")
	res, err := h.service.ItemsByParams(r.Context(), &params)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(res)
}

/*
Query args:

	itemID: always (identifies the exact item)
	category: always (specifies which category to select in)

resp:

	"itemID": int
	"sellerID": string
	"name": string
	"description": string
	"price": int
	"category": string
	"status": string
	"quantity": int
*/
func (h *Handler) ItemByItemID(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	var item structsUFUT.ItemDataRSC
	item.ItemID = q_vals.Get("itemid")
	item.Category = q_vals.Get("category")
	if err := h.service.ItemByItemID(r.Context(), &item); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(item)
}

/*
JSON args:

	"itemID": string (ignored)
	"sellerID": string (ignored)
	"name": string
	"description": string
	"price": int
	"category": string
	"status": string
	"quantity": int

resp:

	"status": "ok"
*/
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item structsUFUT.ItemDataRSC
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	t := item.SellerID
	uid, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	item.ItemID = uid.String()
	item.SellerID = funcsUFUT.GetterIDFromContext(r.Context())
	if err := h.service.CreateItem(r.Context(), &item); err != nil {
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	item.SellerID = t
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

/*
JSON args:

	"itemID": string (always)
	"category": string (always)

resp:

	"status": "ok"
*/
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var item structsUFUT.ItemDataRSC
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteItem(r.Context(), &item); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID     []string `json:"itemID"`
		PassPhrase string   `json:"passphrase"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.PassPhrase != structsUFUT.PASSPHRASE {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var resp struct {
		Successful []bool `json:"successful"`
	}
	resp.Successful = h.service.ReserveItem(r.Context(), req.ItemID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) CancelItemReservation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID     []string `json:"itemID"`
		PassPhrase string   `json:"passphrase"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.PassPhrase != structsUFUT.PASSPHRASE {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err = h.service.CancelItemReservation(r.Context(), req.ItemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
