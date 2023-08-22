package model

import "regexp"

type Redirect struct {
	FromPattern *regexp.Regexp `json:"-"`
	From        string         `json:"from"`
	To          string         `json:"to"`
	Code        int64          `json:"code"`
}

type RedirectApiResponse struct {
	Data struct {
		Redirects []Redirect `json:"redirects"`
	} `json:"data"`
}
