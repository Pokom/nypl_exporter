package nypl_exporter

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	client  *Client
	Queries []string
}

var (
	itemsTotal = prometheus.NewDesc(
		prometheus.BuildFQName("nypl", "items", "total"),
		"Total number of items in the NYPL catalog",
		nil,
		nil,
	)
	collectionsTotal = prometheus.NewDesc(
		prometheus.BuildFQName("nypl", "collections", "total"),
		"Total number of collections in the NYPL catalog",
		[]string{"query", "public_domain"},
		nil,
	)
)

func NewExporter(client *Client) *Exporter {
	return &Exporter{
		client: client,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- itemsTotal
	ch <- collectionsTotal
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Handle /probe requests that pass in queries directly
	if len(e.Queries) > 0 {
		for _, target := range e.Queries {
			total, err := e.client.Search(target, true)
			if err != nil {
				ch <- prometheus.NewInvalidMetric(itemsTotal, err)
				continue
			}
			count, err := strconv.ParseFloat(total.NYPLAPI.Response.NumResults, 64)
			if err != nil {
				ch <- prometheus.NewInvalidMetric(itemsTotal, err)
				continue
			}
			ch <- prometheus.MustNewConstMetric(collectionsTotal, prometheus.GaugeValue, count, target, "true")
		}
		return
	}

	total, err := e.client.ItemsTotal()
	if err != nil {
		log.Printf("Error getting total items: %v", err)
		return
	}
	count, err := strconv.ParseFloat(total.NYPLAPI.Response.Count.Value, 64)
	if err != nil {
		log.Printf("Error parsing total items: %v", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(itemsTotal, prometheus.GaugeValue, count)
}
