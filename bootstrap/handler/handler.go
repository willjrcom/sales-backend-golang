package handler

import "net/http"

type Handler struct {
	Path              string
	Handler           http.Handler
	UnprotectedRoutes []string
}

func NewHandler(path string, handler http.Handler, unprotectedRoutes ...string) *Handler {
	return &Handler{
		Path:              path,
		Handler:           handler,
		UnprotectedRoutes: unprotectedRoutes,
	}
}
