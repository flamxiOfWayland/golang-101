package main

import (
	"log/slog"
	"net/http"
	"slices"

	_ "github.com/flamxiOfWayland/golang-101/api/docs"
	"github.com/flamxiOfWayland/golang-101/api/ping"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

func WithAddressLogs() optFunc {
	return func(c *Controller) error {
		c.router.Use(c.addressLoggingMiddleware)
		return nil
	}
}

func WithAllowedClients(clients []string) optFunc {
	return func(c *Controller) error {
		c.clients = clients
		c.router.Use(c.allowedClientsMiddleware)
		return nil
	}
}

func NewController(opts ...optFunc) *Controller {
	c := &Controller{
		router: mux.NewRouter(),
	}
	c.setup()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type optFunc func(*Controller) error

type Controller struct {
	router  *mux.Router
	clients []string
}

// @title Golang-101 API Gateway
// @version 1.0
// @description This is a sample docs
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.basic BasicAuth

// @contact.name API Support
// @contact.url http://www.wayland.io/support
// @contact.email wayland@wayland.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func (c *Controller) Run() error {
	return http.ListenAndServe("0.0.0.0:80", c.router)
}

func (c *Controller) setup() {
	c.router.HandleFunc("/ping", ping.Handler)
	c.router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
}

func (c *Controller) addressLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Warn("logging middleware", "request uri", r.RequestURI, "request host", r.Host)
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) allowedClientsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(c.clients, r.Host) {
			next.ServeHTTP(w, r)
			slog.Warn("will be serving", "client", r.Host)
			return
		}
		slog.Warn("will not be serving", "client", r.Host)
		r = nil
		return
	})
}
