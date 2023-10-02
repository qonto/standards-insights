package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/internal/build"
	"github.com/qonto/standards-insights/internal/git"
	"github.com/qonto/standards-insights/internal/http"
	"github.com/qonto/standards-insights/internal/metrics"
	"github.com/qonto/standards-insights/internal/providers"
	"github.com/qonto/standards-insights/pkg/checker"
	"github.com/qonto/standards-insights/pkg/daemon"
	"github.com/qonto/standards-insights/pkg/ruler"
	"github.com/spf13/cobra"
)

func serverCmd(configPath, logLevel, logFormat *string) *cobra.Command {
	provider := ""
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			logger := buildLogger(*logLevel, *logFormat)
			config, err := config.New(*configPath)
			exit(err)

			logger.Info(build.VersionMessage())

			registry, ok := prometheus.DefaultRegisterer.(*prometheus.Registry)
			if !ok {
				exit(errors.New("fail to use the Prometheus default registerer"))
			}
			server, err := http.New(registry, logger, config.HTTP)
			exit(err)

			signals := make(chan os.Signal, 1)
			errChan := make(chan error)
			signal.Notify(
				signals,
				syscall.SIGINT,
				syscall.SIGTERM)
			projectMetrics, err := metrics.New(registry, config.Labels)
			exit(err)
			err = server.Start()
			exit(err)

			git, err := git.New(logger, config.Git)
			exit(err)

			ruler := ruler.NewRuler(logger, config.Rules)
			checker := checker.NewChecker(logger, ruler, config.Checks, config.Groups)

			providers, err := providers.NewProviders(logger, config.Providers)
			exit(err)
			if config.Interval == 0 {
				exit(errors.New("the interval configuration is mandatory to run the server"))
			}
			daemon, err := daemon.New(checker, providers, projectMetrics, logger, (time.Duration(config.Interval) * time.Second), git, registry)
			exit(err)
			daemon.Start()

			go func() {
				for sig := range signals {
					switch sig {
					case syscall.SIGINT, syscall.SIGTERM:
						logger.Info(fmt.Sprintf("Received signal %s, shutdown", sig))
						signal.Stop(signals)
						daemon.Stop()
						err := server.Stop()
						if err != nil {
							logger.Error(fmt.Sprintf("Fail to stop: %s", err.Error()))
							errChan <- err
						}
						errChan <- nil
					}
				}
			}()

			exitErr := <-errChan
			exit(exitErr)
			os.Exit(0)
		},
	}

	// TODO: allow to repeat flag
	cmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "Filter providers")

	return cmd
}
