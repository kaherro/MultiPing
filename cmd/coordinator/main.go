package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaherro/MultiPing/internal/checker"
	"github.com/kaherro/MultiPing/internal/config"
	"github.com/kaherro/MultiPing/internal/coordinator/api"
	"github.com/kaherro/MultiPing/internal/coordinator/scheduler"
	"github.com/kaherro/MultiPing/internal/storage/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("loading config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("connect db", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		logger.Error("ping db", "err", err)
		os.Exit(1)
	}

	siteRepo := postgres.NewSiteRepo(pool)
	checkResultRepo := postgres.NewCheckResultRepo(pool)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	sitesHandler := api.NewSitesHandler(siteRepo, logger)
	sitesHandler.Register(mux)
	httpChecker := checker.NewHTTPChecker()
	sched := scheduler.NewScheduler(
		siteRepo,
		checkResultRepo,
		httpChecker,
		8,
		30*time.Second,
		logger,
	)

	go func() {
		if err := sched.Run(ctx); err != nil {
			logger.Error("scheduler stopped with error", "error", err)
			stop()
		}
	}()

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	logger.Info("coordinator started", "addr", cfg.HTTPAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server", "err", err)
		os.Exit(1)
	}
	logger.Info("coordinator stopped")
}
