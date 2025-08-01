package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qonto/standards-insights/internal/providers/aggregates"
	checkeraggregates "github.com/qonto/standards-insights/pkg/checker/aggregates"
	"github.com/qonto/standards-insights/pkg/codeowners"
	"github.com/qonto/standards-insights/pkg/project"
)

type Git interface {
	Clone(url string, ref string, path string) error
	Pull(path string, ref string) error
	SetToken(token string)
}

type Checker interface {
	Run(ctx context.Context, projects []project.Project) []checkeraggregates.ProjectResult
}

type Metrics interface {
	Load(results []checkeraggregates.ProjectResult)
}

type GitProvider interface {
	ConfigureGit(Git) error
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
	// Check if the project is hosted on GitHub
	isGitHubRepo := strings.Contains(project.URL, "github.com")

	if isGitHubRepo {
		d.logger.Debug(fmt.Sprintf("detected GitHub repository: %s", project.URL))
		// Find a GitHub provider to configure Git with the proper token
		foundGitHubProvider := false
		for _, provider := range d.providers {
			if githubProvider, ok := provider.(GitProvider); ok {
				d.logger.Debug(fmt.Sprintf("configuring Git with GitHub token for provider %s", provider.Name()))
				err := githubProvider.ConfigureGit(d.git)
				if err != nil {
					return fmt.Errorf("failed to configure Git with GitHub token: %w", err)
				}
				foundGitHubProvider = true
				break
			}
		}

		if !foundGitHubProvider {
			d.logger.Warn(fmt.Sprintf("no GitHub provider found for repository %s, authentication may fail", project.URL))
		}
	}

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
	registry *prometheus.Registry,
) (*Daemon, error) {
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

func (d *Daemon) tick(configPath string) {
	d.logger.Info("checking projects")
	projects := []project.Project{}
	subProjects := []project.Project{}
	for _, provider := range d.providers {
		providerName := provider.Name()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		providerProjects, err := provider.FetchProjects(ctx)
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			d.logger.Warn(fmt.Sprintf("timeout reached for provider %s: operation took more than 30s", providerName))
		}
		d.logger.Debug(fmt.Sprintf("context canceled for provider %s (fetch status: %v)", providerName, err == nil))
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

			codeowners, err := codeowners.NewCodeowners(proj.Path, configPath)
			if err != nil {
				d.logger.Warn(fmt.Sprintf("Failed to parse CODEOWNERS for project %s: %s", proj.Name, err.Error()))
			}
			// Create subprojects based on expanded paths
			err = d.exploreProjectFiles(proj.Path, codeowners, proj, &subProjects)
			proj.SubProjects = subProjects
			if err != nil {
				d.logger.Warn(fmt.Sprintf("Failed to explore project files for %s: %s", proj.Name, err.Error()))
			}
			projects = append(projects, proj)
		}
	}

	results := d.checker.Run(context.Background(), projects)
	d.metrics.Load(results)
	d.logger.Info("projects checked")
}

func (d *Daemon) exploreProjectFiles(projectPath string, codeowners *codeowners.Codeowners, proj project.Project, projects *[]project.Project) error {
	// Use filepath.Walk to explore all files in the project directory
	return filepath.Walk(projectPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(projectPath, filePath)
			if err != nil {
				return err
			}
			// Assign team owner based on CODEOWNERS
			team, exists := codeowners.GetOwners(relPath)
			if exists {
				labels := make(map[string]string)
				for k, v := range proj.Labels {
					labels[k] = v
				}
				labels["team"] = team

				subProject := project.Project{
					Name:     proj.Name,
					URL:      proj.URL,
					Branch:   proj.Branch,
					Path:     proj.Path,
					FilePath: relPath,
					Labels:   labels,
				}
				*projects = append(*projects, subProject)
			} else {
				// Add the path with the same team as the project team
				labels := make(map[string]string)
				for k, v := range proj.Labels {
					labels[k] = v
				}
				labels["team"] = proj.Labels["team"] // Assuming team is stored in proj.Labels

				subProject := project.Project{
					Name:     proj.Name,
					URL:      proj.URL,
					Branch:   proj.Branch,
					Path:     proj.Path,
					FilePath: relPath,
					Labels:   labels,
				}
				*projects = append(*projects, subProject)
			}
		}
		return nil
	})
}

func (d *Daemon) Start(configPath string) {
	d.logger.Info("starting daemon")
	ticker := time.NewTicker(d.interval)
	d.ticker = ticker
	go func() {
		d.tick(configPath)
		for {
			select {
			case <-d.done:
				return
			case <-ticker.C:
				d.tick(configPath)
			}
		}
	}()
}

func (d *Daemon) Stop() {
	d.logger.Info("stopping daemon")
	d.ticker.Stop()
	d.done <- true
}
