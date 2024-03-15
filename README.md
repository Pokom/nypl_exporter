# NYPL Exporter

New York Public Library exporter is a tool to export metadata from NYPL Digital Collections API as a prometheus metric.
Built as a demo for [Grafana Labs + Friends meetup](https://www.meetup.com/grafana-and-friends-nyc/events/299101735/?utm_medium=referral&utm_campaign=share-btn_savedevents_share_modal&utm_source=link)

## Running

1. First create an account at https://api.repo.nypl.org/ and get an API key.
2. Export the API key as an environment variable `export NYPL_API_KEY=your-api-key`
3. Run the exporter `go run main.go`
4. Visit `http://localhost:8080/metrics` to see the metrics.

## Topics covered in presentation

Over the course of the talk, we covered the following topics:
- What is Prometheus and why it's useful
- What is a Prometheus Exporter
- How to start using the Prometheus go client library 
- Foundational elements for setting up an exporter
  - Prometheus registry
  - Creating metrics and registering them
  - Writing out metrics to the filesystem
  - Creating a simple HTTP server to expose the metrics

## What's next

The final commit of this repository is a working exporter that was built live for the meetup.
It's a solid base to start for Go if you're trying to build your own exporter, however it's not a complete solution.
Here are some ideas to improve it:
- Add operational metrics for the collector(think last duration, error count for api requests, total api requests, etc)
- Create unit tests for the collector
- Create a Grafana Dashboard that can be shared publicly
- Set up a Dockerfile to run the exporter in a container
- Update the http server so that it can gracefully shutdown on SIGTERM/SIGINT signals
