package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qonto/standards-insights/internal/providers/aggregates"
	checkeraggregates "github.com/qonto/standards-insights/pkg/checker/aggregates"
	"github.com/qonto/standards-insights/pkg/project"
)

type Git interface {
	Clone(url string, ref string, path string) error
	Pull(path string, ref string) error
}

type Checker interface {
	Run(ctx context.Context, projects []project.Project) []checkeraggregates.ProjectResult
}

type Metrics interface {
	Load(results []checkeraggregates.ProjectResult)
}

type Daemon struct {
	checker                 Checker
	metrics                 Metrics
	done                    chan bool
	logger                  *slog.Logger
	providers               []aggregates.Provider
	ticker                  *time.Ticker
	interval                time.Duration
	git                     Git
	gitRequestsCounter      *prometheus.CounterVec
	providerRequestsCounter *prometheus.CounterVec
}

func (d *Daemon) cloneOrPull(project project.Project) error {
	_, err := os.Stat(project.Path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		return d.git.Clone(project.URL, project.Branch, project.Path)
	}
	return d.git.Pull(project.Path, project.Branch)
}

func New(checker Checker,
	providers []aggregates.Provider,
	metrics Metrics,
	logger *slog.Logger,
	interval time.Duration,
	git Git,
	registry *prometheus.Registry) (*Daemon, error) {
	gitRequestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "git_requests_total",
			Help: "Number of calls to Git",
		},
		[]string{"status"})
	providerRequestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "provider_requests_total",
			Help: "Number of calls to the projects providers",
		},
		[]string{"name", "status"})
	err := registry.Register(gitRequestsCounter)
	if err != nil {
		return nil, err
	}
	err = registry.Register(providerRequestsCounter)
	if err != nil {
		return nil, err
	}

	return &Daemon{
		done:                    make(chan bool),
		checker:                 checker,
		metrics:                 metrics,
		logger:                  logger,
		providers:               providers,
		interval:                interval,
		git:                     git,
		gitRequestsCounter:      gitRequestsCounter,
		providerRequestsCounter: providerRequestsCounter,
	}, nil
}

func (d *Daemon) tick() {
	d.logger.Info("checking projects")
	projects := []project.Project{}
	for _, provider := range d.providers {
		providerName := provider.Name()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		providerProjects, err := provider.FetchProjects(ctx)
		cancel()
		if err != nil {
			d.logger.Error(fmt.Sprintf("fail to call provider %s: %s", providerName, err.Error()))
			d.providerRequestsCounter.WithLabelValues(providerName, "failure").Inc()
			continue
		}
		d.providerRequestsCounter.WithLabelValues(providerName, "success").Inc()
		for _, proj := range providerProjects {
			err = d.cloneOrPull(proj)
			if err != nil {
				d.logger.Error(fmt.Sprintf("fail to retrieve project %s: %s", proj.Name, err.Error()))
				d.gitRequestsCounter.WithLabelValues("failure").Inc()
				continue
			}
			d.gitRequestsCounter.WithLabelValues("success").Inc()

			labels := make(map[string]string)
			for k, v := range proj.Labels {
				labels[k] = v
			}
			labels["team"] = "bookkeeping"
			subProject := project.Project{
				Name:       proj.Name,
				URL:        proj.URL,
				Branch:     proj.Branch,
				Path:       proj.Path,
				SubProject: "app/announcers/attachment_announcer.rb",
				Labels:     labels,
			}
			projects = append(projects, subProject)

		}
		projects = append(projects, providerProjects...)

	}

	d.logger.Info(fmt.Sprintf("projects: %+v", projects))

	results := d.checker.Run(context.Background(), projects)
	d.metrics.Load(results)
	d.logger.Info("projects checked")
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
