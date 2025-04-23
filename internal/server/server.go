package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func Start(mainRouter *chi.Mux) (*http.Server, *http.Server) {
	config := InitConfig()

	mainServer := createHTTPServer(fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort), mainRouter, config.ServerReadTimeout, config.ServerWriteTimeout)
	go func() {
		log.Info("starting webhook server", zap.String("address", mainServer.Addr))
		if err := mainServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("unable to start webhook server", zap.String("address", mainServer.Addr), zap.Error(err))
		}
	}()

	healthRouter := chi.NewRouter()
	healthRouter.Get("/metrics", promhttp.Handler().ServeHTTP)
	healthRouter.Get("/healthz", HealthCheckHandler)
	healthRouter.Get("/readyz", ReadinessHandler)

	healthServer := createHTTPServer(fmt.Sprintf("%s:%d", config.HealthHost, config.HealthPort), healthRouter, config.ServerReadTimeout, config.ServerWriteTimeout)
	go func() {
		log.Info("starting health server", zap.String("address", healthServer.Addr))
		if err := healthServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("unable to start health server", zap.String("address", healthServer.Addr), zap.Error(err))
		}
	}()

	return mainServer, healthServer
}
