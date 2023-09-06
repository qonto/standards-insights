package http

import (
	"net/http"

	"log/slog"
)

func Health(logger *slog.Logger) func(w http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}
