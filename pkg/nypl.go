package nypl_exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const DefaultURL = "https://api.repo.nypl.org/api/v2/"

type Client struct {
	Key string
	URL string
}

type TotalResponse struct {
	NYPLAPI struct {
		Response struct {
			Count struct {
				Value string `json:"$"`
			} `json:"count"`
		} `json:"response"`
	} `json:"nyplAPI"`
}

func NewClient(key string, url string) *Client {
	return &Client{
		Key: key,
		URL: url,
	}
}

func (c *Client) Do(method string, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.URL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	authorization := fmt.Sprintf("Token token=%q", c.Key)
	req.Header.Add("Authorization", authorization)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return body, nil
}

func (c *Client) ItemsTotal() (float64, error) {
	resp, err := c.Do("GET", "items/total")
	if err != nil {
		return 0.0, err
	}
	var totalResponse *TotalResponse
	if err := json.Unmarshal(resp, &totalResponse); err != nil {
		return 0.0, err
	}
	count, err := strconv.ParseFloat(totalResponse.NYPLAPI.Response.Count.Value, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type SearchResponse struct {
	NYPLAPI struct {
		Response struct {
			NumResults string `json:"numResults"`
			Items      []struct {
				ID     string `json:"uuid"`
				ApiURI string `json:"apiUri"`
				Title  string `json:"title"`
			} `json:"result"`
		} `json:"response"`
	} `json:"nyplAPI"`
}

func (c *Client) Search(q string, publicDomain bool) (float64, error) {
	publicDomainOnly := "true"
	if !publicDomain {
		publicDomainOnly = "false"
	}
	resp, err := c.Do("GET", fmt.Sprintf("items/search?q=%s&publicDomainOnly=%s", q, publicDomainOnly))
	if err != nil {
		return 0.0, err
	}
	var searchResponse *SearchResponse
	if err := json.Unmarshal(resp, &searchResponse); err != nil {
		return 0.0, err
	}
	count, err := strconv.ParseFloat(searchResponse.NYPLAPI.Response.NumResults, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}
