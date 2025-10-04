package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpConfig "github.com/yrss1/workout/internal/pkg/http/config"
	"github.com/yrss1/workout/internal/pkg/log"
)

type Router struct {
	rootRouter chi.Router
	logger     log.Log
}

func NewRouter(
	logger *log.Log,
	config *httpConfig.HTTPConfig,
	recoverer *PanicRecoverer,
) *Router {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		AllowCredentials: config.AllowCredentials,
	})

	rootRouter := chi.NewRouter()
	rootRouter.Use(
		recoverer.Middleware(),
		middleware.RequestID,
		logger.Middleware(),
	)
	rootRouter.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	rootRouter.Route("/api", func(api chi.Router) {
		api.Use(corsMiddleware.Handler)
	})

	var router = &Router{
		logger: *logger,
	}
	return router
}

func (r *Router) Handler() http.Handler {
	return r.rootRouter
}
