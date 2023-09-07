package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qonto/standards-insights/tls"
)

type HTTPServer struct {
	server   *http.Server
	Router   *chi.Mux
	logger   *slog.Logger
	wg       sync.WaitGroup
	registry *prometheus.Registry
}

func New(registry *prometheus.Registry, logger *slog.Logger, config Config) (*HTTPServer, error) {
	var defaultTimeout int = 10
	r := chi.NewRouter()
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if config.WriteTimeout == 0 {
		config.WriteTimeout = defaultTimeout
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultTimeout
	}
	if config.ReadHeaderTimeout == 0 {
		config.ReadHeaderTimeout = defaultTimeout
	}

	server := &http.Server{
		WriteTimeout:      time.Duration(config.WriteTimeout) * time.Second,
		ReadTimeout:       time.Duration(config.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeout) * time.Second,
		Addr:              address,
		Handler:           r,
	}
	if config.KeyPath != "" {
		tlsConfig, err := tls.CreateConfig(config.KeyPath, config.CertPath, config.CacertPath, config.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		logger.Info("enabling HTTPS on the HTTP server")
		if config.ClientAuthType != "" {
			authType, err := tls.GetClientAuthType(config.ClientAuthType)
			if err != nil {
				return nil, err
			}
			tlsConfig.ClientAuth = authType
		}
		server.TLSConfig = tlsConfig
	}

	return &HTTPServer{
		server:   server,
		Router:   r,
		logger:   logger,
		registry: registry,
	}, nil
}

func (h *HTTPServer) Start() error {
	h.logger.Info(fmt.Sprintf("starting HTTP server on %s", h.server.Addr))
	h.Router.Get("/healthz", Health(h.logger))
	h.Router.Method(http.MethodGet, "/metrics", promhttp.Handler())
	go func() {
		defer h.wg.Done()

		var err error
		if h.server.TLSConfig == nil {
			err = h.server.ListenAndServe()
		} else {
			err = h.server.ListenAndServeTLS("", "")
		}

		if err != nil && err != http.ErrServerClosed {
			h.logger.Error(fmt.Sprintf("HTTP server error: %s", err.Error()))
			exitCode := 2
			os.Exit(exitCode)
		}
	}()
	h.wg.Add(1)
	return nil
}

func (h *HTTPServer) Stop() error {
	h.logger.Info("stopping HTTP Server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Error(err.Error())
		return err
	}
	h.wg.Wait()
	return nil
}
