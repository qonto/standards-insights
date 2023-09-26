package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qonto/standards-insights/checks"
	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/daemon"
	"github.com/qonto/standards-insights/git"
	"github.com/qonto/standards-insights/http"
	"github.com/qonto/standards-insights/metrics"
	"github.com/qonto/standards-insights/providers"
	"github.com/qonto/standards-insights/rules"
	"github.com/spf13/cobra"
)

func serverCmd(configPath *string) *cobra.Command {
	provider := ""
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := config.New(*configPath)
			exit(err)
			var programLevel = new(slog.LevelVar)
			programLevel.Set(slog.LevelDebug)
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
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

			ruler := rules.NewRuler(config.Rules)
			checker := checks.NewChecker(ruler, config.Checks, config.Groups)

			providers, err := providers.NewProviders(logger, config.Providers)
			exit(err)
			daemon := daemon.New(checker, providers, projectMetrics, logger, (time.Duration(config.Interval) * time.Second), git)
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
