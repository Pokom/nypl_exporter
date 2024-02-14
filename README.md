# NYPL Exporter

New York Public Library exporter is a tool to export metadata from NYPL Digital Collections API as a prometheus metric.
Built as a demo for Grafana Labs + Friends meetup.

## Running

1. First create an account at https://api.repo.nypl.org/ and get an API key.
2. Export the API key as an environment variable `export NYPL_API_KEY=your-api-key`
3. Run the exporter `go run main.go`
4. Visit `http://localhost:8080/metrics` to see the metrics.
