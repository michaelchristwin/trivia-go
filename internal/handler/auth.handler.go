package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RegisterRoutes() *chi.Mux {
	authRouter := chi.NewRouter()
	authRouter.Post("/login", func(w http.ResponseWriter, r *http.Request) {})
	authRouter.Post("/signup", func(w http.ResponseWriter, r *http.Request) {})

	return authRouter
}
