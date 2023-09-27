package daemon

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/qonto/standards-insights/checks"
	"github.com/qonto/standards-insights/git"
	"github.com/qonto/standards-insights/metrics"
	"github.com/qonto/standards-insights/providers/aggregates"
)

type Daemon struct {
	checker        *checks.Checker
	projectMetrics *metrics.Project
	done           chan bool
	logger         *slog.Logger
	providers      []aggregates.Provider
	ticker         *time.Ticker
	interval       time.Duration
	git            *git.Git
}

func (d *Daemon) cloneOrPull(project aggregates.Project) error {
	_, err := os.Stat(project.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		return d.git.Clone(project.URL, project.Branch, project.Path)
	}
	return d.git.Pull(project.Path, project.Branch)
}

func New(checker *checks.Checker, providers []aggregates.Provider, projectMetrics *metrics.Project, logger *slog.Logger, interval time.Duration, git *git.Git) *Daemon {
	return &Daemon{
		done:           make(chan bool),
		checker:        checker,
		projectMetrics: projectMetrics,
		logger:         logger,
		providers:      providers,
		interval:       interval,
		git:            git,
	}
}

func (d *Daemon) tick() {
	d.logger.Debug("computing projects metrics")
	var err error
	projects := []aggregates.Project{}
	for _, provider := range d.providers {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		providerProjects, err := provider.FetchProjects(ctx)
		cancel()
		// TODO clean + alert
		if err != nil {
			d.logger.Error(err.Error())
			break
		}
		for _, project := range providerProjects {
			err = d.cloneOrPull(project)
			if err != nil {
				d.logger.Error(err.Error())
				break
			}
		}
		projects = append(projects, providerProjects...)
	}
	if err == nil {
		results := d.checker.Run(context.Background(), projects)
		d.projectMetrics.LoadProjectsMetrics(results)
		d.logger.Debug("metrics computed")
	}
}

func (d *Daemon) Start() {
	d.logger.Info("starting daemon")
	ticker := time.NewTicker(d.interval)
	d.ticker = ticker
	go func() {
		d.tick()
		for {
			select {
			case <-d.done:
				return
			case <-ticker.C:
				d.tick()
			}
		}
	}()
}

func (d *Daemon) Stop() {
	d.logger.Info("stopping daemon")
	d.ticker.Stop()
	d.done <- true
}
