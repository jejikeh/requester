package internal

import (
	"net/http"
	"strings"
)

func NewHttpClient() *http.Client {
	return &http.Client{}
}

func MakeRequest(t *Task) (*http.Response, error) {
	client := NewHttpClient()

	req, err := http.NewRequest(t.Method, t.URL, strings.NewReader(t.Body))

	for k, v := range t.RequestHeaders {
		req.Header.Set(k, v)
	}

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
