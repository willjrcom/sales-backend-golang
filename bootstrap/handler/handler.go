package handler

import "net/http"

type Handler struct {
	Path    string
	Handler http.Handler
}

func NewHandler(path string, handler http.Handler) *Handler {
	return &Handler{
		Path:    path,
		Handler: handler,
	}
}
