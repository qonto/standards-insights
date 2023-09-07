package daemon

import (
	"context"
	"log/slog"
	"time"

	"github.com/qonto/standards-insights/checks"
	"github.com/qonto/standards-insights/metrics"
	"github.com/qonto/standards-insights/providers/aggregates"
)

type Daemon struct {
	checker        *checks.Checker
	projectMetrics *metrics.Project
	done           chan bool
	logger         *slog.Logger
	projects       []aggregates.Project
	ticker         *time.Ticker
	interval       time.Duration
}

func New(checker *checks.Checker, projects []aggregates.Project, projectMetrics *metrics.Project, logger *slog.Logger, interval time.Duration) *Daemon {
	return &Daemon{
		done:           make(chan bool),
		checker:        checker,
		projectMetrics: projectMetrics,
		logger:         logger,
		projects:       projects,
		interval:       interval,
	}
}

func (d *Daemon) Start() {
	d.logger.Info("starting daemon")
	ticker := time.NewTicker(d.interval)
	d.ticker = ticker
	go func() {
		for {
			select {
			case <-d.done:
				return
			case <-ticker.C:
				d.logger.Info("computing projects metrics")
				results := d.checker.Run(context.Background(), d.projects)
				d.projectMetrics.LoadProjectsMetrics(results)
			}
		}
	}()
}

func (d *Daemon) Stop() {
	d.logger.Info("stopping daemon")
	d.ticker.Stop()
	d.done <- true
}
