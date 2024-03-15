package nypl_exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	NYPLTotalItemsDesc = prometheus.NewDesc(
		"nypl_items_total",
		"The total number of items in the NYPL",
		nil, nil)
	NYPLCollectionsTotalDesc = prometheus.NewDesc(
		"nypl_collections_total",
		"The total number of items in the NYPL",
		[]string{"collection"}, nil)
)

type NYPLCollector struct {
	Client *Client
}

func (N NYPLCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- NYPLTotalItemsDesc
	ch <- NYPLCollectionsTotalDesc
}

func (N NYPLCollector) Collect(ch chan<- prometheus.Metric) {
	total, err := N.Client.ItemsTotal()
	if err != nil {
		fmt.Println(err)
	}
	ch <- prometheus.MustNewConstMetric(
		NYPLTotalItemsDesc,
		prometheus.GaugeValue,
		total)

	collections := []string{"cats", "dogs", "birds", "basil"}
	for _, collection := range collections {
		total, err := N.Client.Search(collection, true)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			NYPLCollectionsTotalDesc,
			prometheus.GaugeValue,
			total,
			collection)
	}
}
