package httpserver

import (
	"net/http"
)

type Router interface {
	Default()
	ServeHTTP()
	AddPrefix(prefix string) Router
	AllowRecovery() Router
	AllowLog() Router
	AllowHealthCheck() Router
	AllowCors() Router
	AddPath(path, method string, handler http.HandlerFunc) Router
	AddMiddleware(middleware func(next http.Handler) http.Handler) Router
}
