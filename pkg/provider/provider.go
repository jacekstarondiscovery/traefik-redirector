package provider

import (
	"bytes"
	"encoding/json"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/client"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/model"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Provider struct {
	log    *log.Logger
	client client.HTTPClient
}

func New(log *log.Logger, client client.HTTPClient) *Provider {
	return &Provider{
		log:    log,
		client: client,
	}
}

func (p *Provider) GetRedirects(method, endpoint, data string) ([]model.Redirect, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	redResp := &model.RedirectApiResponse{}
	err = json.Unmarshal(body, redResp)
	if err != nil {
		return nil, err
	}

	return p.build(redResp), nil
}

func (p *Provider) build(redResp *model.RedirectApiResponse) []model.Redirect {
	result := []model.Redirect{}
	for _, redir := range redResp.Data.Redirects {
		fromPattern, err := regexp.Compile(redir.From)
		if err != nil {
			p.log.Println("Unable to compile regexp: ", redir.From)
			continue
		}

		result = append(result, model.Redirect{
			FromPattern: fromPattern,
			From:        redir.From,
			To:          redir.To,
			Code:        redir.Code,
		})
	}

	return result
}
