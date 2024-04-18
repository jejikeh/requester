package internal

import (
	"log"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) MakeRequest(t *Task) (*http.Response, error) {
	log.Printf("Requesting %s %s", t.Method, t.URL)

	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(t.Body))

	for k, v := range t.RequestHeaders {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
