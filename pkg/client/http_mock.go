package client

import (
	"bytes"
	"io"
	"net/http"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	responseBody := io.NopCloser(bytes.NewReader([]byte("{\"data\":{\"redirects\":[{\"from\":\"^/xxx\",\"to\":\"https://tvn24.pl\",\"code\":307}]}}")))
	return &http.Response{
		StatusCode: 200,
		Body:       responseBody,
	}, nil
}
