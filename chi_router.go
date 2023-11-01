package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ChiHandle struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

type ChiRouter struct {
	port        string
	router      *chi.Mux
	prefix      string
	middleware  []func(next http.Handler) http.Handler
	handlers    map[string]ChiHandle
	healthCheck bool
	logRequest  bool
	cors        bool
	recovery    bool
}

func NewChiRouter(port string) Router {
	return &ChiRouter{
		port:       port,
		router:     chi.NewRouter(),
		middleware: make([]func(next http.Handler) http.Handler, 0),
		handlers:   make(map[string]ChiHandle, 0),
	}
}

func (r *ChiRouter) Default() {
	r.
		AllowCors().
		AllowHealthCheck().
		AllowLog().
		AllowRecovery().
		ServeHTTP()
}

func (r *ChiRouter) AddPrefix(prefix string) Router {
	r.prefix = prefix
	return r
}

func (r *ChiRouter) AllowRecovery() Router {
	r.recovery = true
	return r
}

func (r *ChiRouter) AllowLog() Router {
	r.logRequest = true
	return r
}

func (r *ChiRouter) AllowHealthCheck() Router {
	r.healthCheck = true
	return r
}

func (r *ChiRouter) AllowCors() Router {
	r.cors = true
	return r
}

func (r *ChiRouter) ServeHTTP() {

	// middleware
	if r.logRequest {
		r.AddMiddleware(middleware.Logger)
	}
	if r.cors {
		r.AddMiddleware(r.accessControlMiddleware)
	}
	if r.recovery {
		r.AddMiddleware(middleware.Recoverer)
	}

	for _, h := range r.middleware {
		r.router.Use(h)
	}

	// handler
	if r.healthCheck {
		r.AddPath("/health", "GET", r.healthCheckHandler)
	}

	if r.prefix != "" {
		log.Printf("version-api %s", r.prefix)
	}
	r.router.Route(r.prefix, func(router chi.Router) {
		for _, h := range r.handlers {
			log.Printf("api: %s, method: %s", h.Path, h.Method)
			router.Method(h.Method, h.Path, h.Handler)
		}
	})

	// server
	server := http.Server{
		Addr:         ":" + r.port,
		Handler:      r.router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Println("Server started on: " + r.port)
		errs <- server.ListenAndServe()
	}()

	log.Println("exit", <-errs)
}

func (ChiRouter) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ChiJson(w, 200, map[string]interface{}{
		"message": "service is running",
	})
}

func (ChiRouter) accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (r *ChiRouter) AddPath(path, method string, handler http.HandlerFunc) Router {
	r.handlers[path+"_"+method] = ChiHandle{
		Path:    path,
		Method:  method,
		Handler: handler,
	}
	return r
}

func (r *ChiRouter) AddMiddleware(middleware func(next http.Handler) http.Handler) Router {
	r.middleware = append(r.middleware, middleware)
	return r
}

func ChiJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
