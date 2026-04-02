package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/SKIND0A/online-shop/internal/repository/postgres"
	"github.com/SKIND0A/online-shop/internal/usecase"
)

type AuthHandler struct {
	auth *usecase.AuthUsecase
}

func NewAuthHandler(auth *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	//начало
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var in usecase.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	res, err := h.auth.Register(r.Context(), in)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request")
			return
		}
		if errors.Is(err, postgres.ErrEmailAlreadyExists) {
			writeError(w, http.StatusConflict, "CONFLICT", "email already exists")
			return
		}
		log.Printf("auth register: %v", err)
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeSuccess(w, http.StatusCreated, res)
	//конец
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	//начало
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var in usecase.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request body")
		return
	}

	res, err := h.auth.Login(r.Context(), in)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid request")
			return
		}
		if errors.Is(err, usecase.ErrInvalidCredentials) || errors.Is(err, usecase.ErrInactiveUser) {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeSuccess(w, http.StatusOK, res)
	//конец
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	raw := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(raw, prefix) {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token")
		return
	}
	token := strings.TrimSpace(raw[len(prefix):])
	if token == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing bearer token")
		return
	}

	res, err := h.auth.Me(r.Context(), token)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidToken) {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeSuccess(w, http.StatusOK, res)
}

func writeSuccess(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    data,
		"error":   nil,
	})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"data":    nil,
		"error": map[string]any{
			"code":    code,
			"message": message,
			"details": nil,
		},
	})
}
