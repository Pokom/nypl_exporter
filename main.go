package main

import (
	"fmt"
	"os"

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

func run(client *nypl_exporter.Client) error {
	resp, err := client.ItemsTotal()
	if err != nil {
		return err
	}
	fmt.Println(resp.NYPLAPI.Response.Count.Value)
	return nil
}
