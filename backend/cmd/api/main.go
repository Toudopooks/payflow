package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Toudopooks/payflow/backend/internal/config"
	appmiddleware "github.com/Toudopooks/payflow/backend/internal/middleware"
	"github.com/Toudopooks/payflow/backend/internal/observability"
	mysqlstore "github.com/Toudopooks/payflow/backend/internal/repository/mysql"
	"github.com/Toudopooks/payflow/backend/internal/service"
	"github.com/Toudopooks/payflow/backend/internal/transport/httpapi"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "healthcheck" {
		response, err := http.Get("http://127.0.0.1:8080/health")
		if err != nil || response.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		_ = response.Body.Close()
		return
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		log.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	store, err := mysqlstore.Open(cfg.DatabaseDSN)
	if err != nil {
		log.Error("database setup failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	metrics := observability.New(prometheus.DefaultRegisterer)
	svc := service.New(store)
	handler := httpapi.New(svc, metrics, cfg.DemoMode, cfg.AllowedOrigins, appmiddleware.Observe(log, metrics))
	server := &http.Server{Addr: ":" + cfg.Port, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	go func() {
		log.Info("api listening", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdown); err != nil {
		log.Error("graceful shutdown failed", "error", err)
	}
}
