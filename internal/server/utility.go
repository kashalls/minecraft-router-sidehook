package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kashalls/minecraft-router-sidehook/internal/log"
	"go.uber.org/zap"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func createHTTPServer(addr string, hand http.Handler, readTimeout, writeTimeout time.Duration) *http.Server {
	return &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         addr,
		Handler:      hand,
	}
}

// ShutdownGracefully gracefully shutdown the http server
func ShutdownGracefully(mainServer *http.Server, healthServer *http.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigCh

	log.Info("shutting down servers due to received signal", zap.Any("signal", sig))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := mainServer.Shutdown(ctx); err != nil {
		log.Error("error shutting down main server", zap.Error(err))
	}

	if err := healthServer.Shutdown(ctx); err != nil {
		log.Error("error shutting down health server", zap.Error(err))
	}
}
