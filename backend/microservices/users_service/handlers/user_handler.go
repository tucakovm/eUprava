package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"users_module/models"
	"users_module/services"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	Svc services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{Svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json")
		return
	}
	user, err := h.Svc.Register(r.Context(), req) // prosleđujemo ctx
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserExists):
			httpError(w, http.StatusConflict, "email or username already exists")
		default:
			httpError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid json")
		return
	}
	resp, err := h.Svc.Login(r.Context(), req) // prosleđujemo ctx
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			httpError(w, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, services.ErrUserDisabled):
			httpError(w, http.StatusForbidden, "user disabled")
		default:
			httpError(w, http.StatusBadRequest, err.Error())
		}
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (a *AuthHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := a.Svc.GetUser(r.Context(), id)

	if err != nil {
		http.Error(rw, "Database exception", http.StatusInternalServerError)
		return
	}

	writeJSON(rw, http.StatusOK, user)
}

// helpers

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
