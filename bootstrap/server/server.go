package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
)

type ServerInterface interface {
	newServer()
	StartServer(port string)
	AddHandler(handler *handler.Handler)
}

type ServerChi struct {
	Router     *chi.Mux
	HttpServer *http.Server
}

func NewServerChi() *ServerChi {
	c := &ServerChi{}
	c.newServer()
	return c
}

func (c *ServerChi) newServer() {
	c.Router = chi.NewRouter()
}

func (c *ServerChi) StartServer(port string) error {
	// create http server with handler from router
	c.HttpServer = &http.Server{
		Addr:              port,
		Handler:           c.Router,
		ReadHeaderTimeout: 30 * time.Second,
	}

	if err := c.HttpServer.ListenAndServe(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (c *ServerChi) AddHandler(h *handler.Handler) {
	c.Router.Mount(h.Path, h.Handler)
}