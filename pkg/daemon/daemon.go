package daemon

import (
	"bufio"
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

func (d *Daemon) parseCodeowners(projectPath string) (map[string]string, error) {
	codeownersPath := filepath.Join(projectPath, ".gitlab", "CODEOWNERS")
	file, err := os.Open(codeownersPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty map if CODEOWNERS does not exist
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer file.Close()

	pathOwners := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			path := parts[0]
			team := strings.TrimPrefix(parts[1], "@")
			
			// Add the original path-team mapping
			if _, exists := pathOwners[path]; !exists {
				pathOwners[path] = team
			}
			
			// Expand paths and add them to pathOwners
			err := d.expandPaths(projectPath, path, team, pathOwners)
			if err != nil {
				d.logger.Warn(fmt.Sprintf("Failed to expand path %s: %s", path, err.Error()))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return pathOwners, nil
}

func (d *Daemon) expandPaths(projectPath, pattern string, team string, pathOwners map[string]string) error {
	// Use Glob to find all matches for the pattern
	matches, err := filepath.Glob(filepath.Join(projectPath, pattern))
	if err != nil {
		return err
	}
	for _, match := range matches {
		if info, err := os.Stat(match); err == nil && info.IsDir() {
			// Walk through the directory to find all files
			err := filepath.Walk(match, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					relPath, err := filepath.Rel(projectPath, filePath)
					if err != nil {
						return err
					}
					if _, exists := pathOwners[relPath]; !exists {
						pathOwners[relPath] = team
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else if info != nil && !info.IsDir() {
			// If it's a file, add it directly
			relPath, err := filepath.Rel(projectPath, match)
			if err != nil {
				return err
			}
			if _, exists := pathOwners[relPath]; !exists {
				pathOwners[relPath] = team
			}
		}
	}
	return nil
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
	subProjects := []project.Project{}
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

			codeownersLabels, err := d.parseCodeowners(proj.Path)
			if err != nil {
				d.logger.Warn(fmt.Sprintf("Failed to parse CODEOWNERS for project %s: %s", proj.Name, err.Error()))
			}
			
			// Create subprojects based on expanded paths
			err = d.exploreProjectFiles(proj.Path, codeownersLabels, proj, &subProjects)
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

func (d *Daemon) exploreProjectFiles(projectPath string, codeownersLabels map[string]string, proj project.Project, projects *[]project.Project) error {
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
			team, exists := codeownersLabels[relPath]
			if exists {
				labels := make(map[string]string)
				for k, v := range proj.Labels {
					labels[k] = v
				}
				labels["team"] = team

				subProject := project.Project{
					Name:       proj.Name,
					URL:        proj.URL,
					Branch:     proj.Branch,
					Path:       proj.Path,
					SubProject: relPath,
					Labels:     labels,
				}
				*projects = append(*projects, subProject)
			} else {
				// Add the path as a subproject with the same team as the project team
				labels := make(map[string]string)
				for k, v := range proj.Labels {
					labels[k] = v
				}
				labels["team"] = proj.Labels["team"] // Assuming team is stored in proj.Labels

				subProject := project.Project{
					Name:       proj.Name,
					URL:        proj.URL,
					Branch:     proj.Branch,
					Path:       proj.Path,
					SubProject: relPath,
					Labels:     labels,
				}
				*projects = append(*projects, subProject)
			}
		}
		return nil
	})
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