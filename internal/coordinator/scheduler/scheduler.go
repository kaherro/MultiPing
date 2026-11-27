package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kaherro/MultiPing/internal/checker"
	"github.com/kaherro/MultiPing/internal/domain"
	"github.com/kaherro/MultiPing/internal/probe/runner"
	"github.com/kaherro/MultiPing/internal/storage/postgres"
)

type Scheduler struct {
	sites *postgres.SiteRepo
	results *postgres.CheckResultRepo
	checker checker.Checker	
	workers int
	interval time.Duration
	logger *slog.Logger
}

func NewScheduler (sites *postgres.SiteRepo, results *postgres.CheckResultRepo, 
	checker checker.Checker, workers int, interval time.Duration, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		sites: sites, 
		results: results, 
		checker: checker, 
		workers: workers, 
		interval: interval,
		logger: logger,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	if s.workers < 1 {
		return fmt.Errorf("workers must be greater than zero")
	}
	if s.interval <= 0 {
		return fmt.Errorf("interval must be greater than zero")
	}
	if err := s.runOnce(ctx); err != nil && ctx.Err() == nil {
		s.logger.Error("scheduler failed to run probe", "error", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped")
			return nil
		case <-ticker.C:
			if err := s.runOnce(ctx); err != nil && ctx.Err() == nil {
				s.logger.Error("scheduler failed to run probe", "error", err)
			}
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) error {
	sites, err := s.sites.List(ctx)
	if err != nil {
		return fmt.Errorf("load sites: %w", err)
	}
	if len(sites) == 0 {
		return nil
	}

	jobs := make(chan domain.Site, len(sites))
	results := make(chan domain.CheckResult, len(sites))

	for _, site := range sites {
		jobs <- site
	}
	close(jobs)
	
	go runner.RunPool(ctx, s.checker, s.workers, jobs, results)

	for result := range results {
		if err := s.results.Create(ctx, &result); err != nil {
			return fmt.Errorf("saved check result for site %d: %w", result.SiteID, err)
		}
	}

	return nil
}