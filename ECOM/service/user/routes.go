package user

import (
	"crypto/ecdh"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleLogin).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, *http.Request) {

}