package auth

import (
	"encoding/json"
	"net/http"
	structsUFUT "ufut/lib/structs"
)

type Handler struct {
	service *Service
}

func NewHandler(srvc *Service) *Handler {
	return &Handler{service: srvc}
}

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	handledFuncs := map[string]http.HandlerFunc{
		"POST /api/user/authUser":         h.AuthUser,
		"POST /api/user/registerUser":     h.RegisterUser,
		"POST /api/user/updateUserPasswd": h.UpdateUserPasswd,
		"POST /api/user/updateJWTUser":    h.UpdateJWTUser,

		"POST /api/staff/authStaff":         h.AuthStaff,
		"POST /api/staff/registerStaff":     h.RegisterStaff,
		"POST /api/staff/updateStaffPasswd": h.UpdateStaffPasswd,
		"POST /api/user/updateJWTStaff":     h.UpdateJWTStaff,
	}

	for key, val := range handledFuncs {
		mux.HandleFunc(key, val)
	}
}

/*
JSON args:

	"login": string
	"password": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
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
	resp, err := h.service.AuthUser(r.Context(), &req)
	if err != nil {
		http.Error(w, "authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"login": string
	"password": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
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
	resp, err := h.service.AuthStaff(r.Context(), &req)
	if err != nil {
		http.Error(w, "authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"login": string
	"password": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
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
	resp, err := h.service.RegisterUser(r.Context(), &req)
	if err != nil {
		http.Error(w, "registration failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"login": string
	"password": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
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
	resp, err := h.service.RegisterStaff(r.Context(), &req)
	if err != nil {
		http.Error(w, "registration failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"login": string			/jwt
	"password": string
	"newPassword": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
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
	resp, err := h.service.UpdateUserPasswd(r.Context(), &req)
	if err != nil {
		http.Error(w, "password update failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"login": string			/jwt
	"password": string
	"newPassword": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
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
	resp, err := h.service.UpdateStaffPasswd(r.Context(), &req)
	if err != nil {
		http.Error(w, "password update failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

/*
JSON args:

	"rt": string

resp:

	"jwt": string (JSON Web Token)
	"rt": string (Refresh Token)
*/
func (h *Handler) UpdateJWTUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RT string `json:"rt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	res, err := h.service.UpdateJWTUser(r.Context(), req.RT)
	if err != nil {
		http.Error(w, "Invalid RT", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(res)
}

func (h *Handler) UpdateJWTStaff(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RT string `json:"rt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	res, err := h.service.UpdateJWTStaff(r.Context(), req.RT)
	if err != nil {
		http.Error(w, "Invalid RT", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(res)
}
