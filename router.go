package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	Default()
	ServeHTTP()
	GetRouter() *chi.Mux
	AddPrefix(prefix string) Router
	AllowRecovery() Router
	AllowLog() Router
	AllowHealthCheck() Router
	AllowCors() Router
	AddPath(path, method string, handler http.HandlerFunc) Router
	AddMiddleware(middleware func(next http.Handler) http.Handler) Router
}
