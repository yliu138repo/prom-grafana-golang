package httpserver

import (
	"fmt"
	"goprometheus/internal/instrumentation"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func configureRouter(userHandler *UserHandler) *chi.Mux {

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.Timeout(500 * time.Second))
	r.Use(middleware.RealIP)
	r.Use(instrumentation.PrometheusMiddleware)
	r.Use(middleware.Heartbeat("/ping"))

	r.Handle("/metrics", promhttp.Handler())

	//r.NotFound(notFoundHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Matrix !"))
	})

	// Add routes
	r.Get("/api/v1/order", getOrderHandler)
	r.Get("/api/v1/user/{id}", userHandler.GetUser)
	r.Post("/api/v1/user", userHandler.CreateUser)
	r.Get("/api/v1/users", userHandler.GetAllUsers)

	return r
}

func NewHTTPServer(userHandler *UserHandler, port int) *http.Server {
	router := configureRouter(userHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
