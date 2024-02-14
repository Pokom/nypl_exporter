package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"nypl_exporter/pkg"
)

const (
	defaultUrl = "https://api.repo.nypl.org/api/v2/"
)

func main() {
	key := os.Getenv("NYPL_API_KEY")
	if key == "" {
		fmt.Println("NYPL_API_KEY must be set")
		os.Exit(1)
	}
	url := os.Getenv("NYPL_API_URL")
	if url == "" {
		fmt.Printf("NYPL_API_URL not set, defaulting to %q\n", defaultUrl)
		url = defaultUrl
	}
	client := nypl_exporter.NewClient(key, url)
	if err := run(client); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleProbe(w http.ResponseWriter, r *http.Request, client *nypl_exporter.Client) {
	q := r.URL.Query()
	query := q["query"]
	if len(query) == 0 {
		w.WriteHeader(http.StatusBadRequest)
	}
	var queries []string
	for _, qm := range query {
		// Potentially split on comma to allow multiple queries in a single entry
		queries = append(queries, strings.Split(qm, ",")...)
	}

	collector := nypl_exporter.NewExporter(client)
	collector.Queries = queries
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func run(client *nypl_exporter.Client) error {
	exp := nypl_exporter.NewExporter(client)
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		exp,
	)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.Handle("/probe", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleProbe(w, r, client)
	}))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}
