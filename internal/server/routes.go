package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initRoutes(metricsPath string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(metricsPath, promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Prometheus Directory Size Exporter is up and running"))
	})

	return mux
}
