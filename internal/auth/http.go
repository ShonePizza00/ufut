package auth

import (
	"context"
	"encoding/json"
	"net/http"
	structsUFUT "ufut/lib"
)

type Handler struct {
	service *Service
	ctx     context.Context
}

func NewHandler(ctx context.Context, srvc *Service) *Handler {
	return &Handler{service: srvc, ctx: ctx}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /api/authUser", h.AuthUser)
	mux.HandleFunc("POST /api/authStaff", h.AuthStaff)
	mux.HandleFunc("POST /api/registerUser", h.RegisterUser)
	mux.HandleFunc("POST /api/registerStaff", h.RegisterStaff)
	mux.HandleFunc("POST /api/updateUserPasswd", h.UpdateUserPasswd)
	mux.HandleFunc("POST /api/updateStaffPasswd", h.UpdateStaffPasswd)
	mux.HandleFunc("GET /api/verifyTokenUser", h.VerifyTokenUser)
	mux.HandleFunc("GET /api/verifyTokenStaff", h.VerifyTokenStaff)
}

func (h *Handler) AuthUser(w http.ResponseWriter, r *http.Request) {

	var req structsUFUT.UpdatePasswdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Passwd == "" {
		http.Error(w, "missing login and password", http.StatusBadRequest)
		return
	}
	token, err := h.service.AuthUser(h.ctx, &req)
	if err != nil {
		http.Error(w, "authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(structsUFUT.TokenResponse{Token: token})
}

func (h *Handler) AuthStaff(w http.ResponseWriter, r *http.Request) {
	var req structsUFUT.UpdatePasswdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Passwd == "" {
		http.Error(w, "missing login and password", http.StatusBadRequest)
		return
	}
	token, err := h.service.AuthStaff(h.ctx, &req)
	if err != nil {
		http.Error(w, "authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(structsUFUT.TokenResponse{Token: token})
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req structsUFUT.UpdatePasswdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Passwd == "" {
		http.Error(w, "missing login and password", http.StatusBadRequest)
		return
	}
	token, err := h.service.RegisterUser(h.ctx, &req)
	if err != nil {
		http.Error(w, "registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(structsUFUT.TokenResponse{Token: token})
}

func (h *Handler) RegisterStaff(w http.ResponseWriter, r *http.Request) {
	var req structsUFUT.UpdatePasswdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Passwd == "" {
		http.Error(w, "missing login and password", http.StatusBadRequest)
		return
	}
	token, err := h.service.RegisterStaff(h.ctx, &req)
	if err != nil {
		http.Error(w, "registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(structsUFUT.TokenResponse{Token: token})
}

func (h *Handler) UpdateUserPasswd(w http.ResponseWriter, r *http.Request) {
	var req structsUFUT.UserUpdatePasswd
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Passwd == "" || req.NewPasswd == "" {
		http.Error(w, "missing login, password or new password", http.StatusBadRequest)
		return
	}
	token, err := h.service.UpdateUserPasswd(h.ctx, &req)
	if err != nil {
		http.Error(w, "password update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(structsUFUT.TokenResponse{Token: token})
}

func (h *Handler) UpdateStaffPasswd(w http.ResponseWriter, r *http.Request) {
	var req structsUFUT.UserUpdatePasswd
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Passwd == "" || req.NewPasswd == "" {
		http.Error(w, "missing login, password or new password", http.StatusBadRequest)
		return
	}
	token, err := h.service.UpdateStaffPasswd(h.ctx, &req)
	if err != nil {
		http.Error(w, "password update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(structsUFUT.TokenResponse{Token: token})
}

func (h *Handler) VerifyTokenUser(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	token := q_vals.Get("token")
	passphrase := q_vals.Get("passphrase")
	if passphrase != structsUFUT.PASSPHRASE {
		http.Error(w, "Access Denied!", http.StatusForbidden)
		return
	}
	if token == "" {
		http.Error(w, "missing authorization token", http.StatusBadRequest)
		return
	}
	valid, err := h.service.VerifyTokenUser(h.ctx, token)
	if err != nil {
		http.Error(w, "token verification failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"userID": valid})
}

func (h *Handler) VerifyTokenStaff(w http.ResponseWriter, r *http.Request) {
	q_vals := r.URL.Query()
	token := q_vals.Get("token")
	passphrase := q_vals.Get("passphrase")
	if passphrase != structsUFUT.PASSPHRASE {
		http.Error(w, "Access Denied!", http.StatusForbidden)
		return
	}
	if token == "" {
		http.Error(w, "missing authorization token", http.StatusBadRequest)
		return
	}
	valid, err := h.service.VerifyTokenStaff(h.ctx, token)
	if err != nil {
		http.Error(w, "token verification failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"userID": valid})
}
