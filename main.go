package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	nypl_exporter "nypl_exporter/pkg"
)

var (
	MetricsRequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "metrics_request_total",
			Help: "The total number of metrics requests",
		})
	MetricsRequestDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "metrics_request_duration_seconds",
			Help: "The duration of metrics requests",
		})
	MetricsRequestDurationHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "metrics_request_duration_histogram_seconds",
			Help:    "The duration of metrics requests",
			Buckets: prometheus.DefBuckets,
		})
	MetricsRequestSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "metrics_request_summary_seconds",
			Help:       "The duration of metrics requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
)

func main() {
	apiKey := os.Getenv("NYPL_API_KEY")
	if apiKey == "" {
		fmt.Println("NYPL_API_KEY environment variable is required")
		os.Exit(1)
	}
	if err := run(apiKey); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(apiKey string) error {
	collector := &nypl_exporter.NYPLCollector{
		Client: nypl_exporter.NewClient(apiKey, nypl_exporter.DefaultURL),
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(MetricsRequestCounter,
		MetricsRequestDuration,
		MetricsRequestDurationHistogram,
		collector,
	)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		MetricsRequestCounter.Inc()
		handler.ServeHTTP(w, r)
		MetricsRequestDuration.Set(time.Since(start).Seconds())
		MetricsRequestDurationHistogram.Observe(time.Since(start).Seconds())
		MetricsRequestSummary.Observe(time.Since(start).Seconds())
	})
	return http.ListenAndServe(":8080", nil)
}
