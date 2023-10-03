package traefik_redirector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/cache"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/client"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/log"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/model"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/provider"
	"net/http"
	"time"
)

type Config struct {
	MaxAge             int64     `json:"maxAge"`
	CacheControlMaxAge int       `json:"cacheControlMaxAge"`
	Endpoint           string    `json:"endpoint"`
	Method             string    `json:"method"`
	Data               string    `json:"data,omitempty"`
	ClientType         string    `json:"clientType"`
	DebugParameter     string    `json:"debugParameter"`
	LogLevel           log.Level `json:"logLevel"`
}

func CreateConfig() *Config {
	return &Config{
		MaxAge:             60,
		CacheControlMaxAge: 60,
		ClientType:         "mock",
		Method:             "GET",
		Endpoint:           "",
		Data:               "",
		DebugParameter:     "debug",
		LogLevel:           log.Debug,
	}
}

type TraefikRedirector struct {
	log      *log.Logger
	config   *Config
	next     http.Handler
	name     string
	cache    *cache.RedirectCache
	provider *provider.Provider
}

func (r *TraefikRedirector) UpdateRedirects() bool {
	r.log.Debug("[UpdateRedirects] Update redirects")
	if !r.cache.Lock() {
		r.log.Error("[UpdateRedirects] Cache is locked")
		return false
	}

	defer r.cache.Unlock()

	redirects, err := r.provider.GetRedirects(r.config.Method, r.config.Endpoint, r.config.Data)
	if err != nil {
		r.log.Error("[UpdateRedirects] Fetch error", fmt.Sprintf("%+v", redirects))
		return false
	}

	r.log.Debug("[UpdateRedirects] Update Cache with: ", fmt.Sprintf("%+v", redirects))
	r.cache.Update(redirects)

	return true
}

func (r *TraefikRedirector) GetRedirects() []model.Redirect {
	if r.cache.IsFresh() {
		return r.cache.Redirects
	}

	go r.UpdateRedirects()

	return r.cache.Redirects
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	logger := log.New(config.LogLevel)
	logger.Debug("[New] with: ", fmt.Sprintf("%+v", config))

	var cl client.HTTPClient
	if config.ClientType == "mock" {
		cl = &client.MockClient{}
	} else {
		cl = &http.Client{}
	}

	tr := &TraefikRedirector{
		log:      logger,
		config:   config,
		next:     next,
		name:     name,
		cache:    cache.New(config.MaxAge, time.Now),
		provider: provider.New(logger, cl),
	}

	return tr, nil
}

func (r *TraefikRedirector) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.ServeLoad(rw, req) {
		return
	}

	if r.ServeDebug(rw, req) {
		return
	}

	if r.ServeRedirect(rw, req) {
		return
	}

	r.next.ServeHTTP(rw, req)
}

func (r *TraefikRedirector) ServeRedirect(rw http.ResponseWriter, req *http.Request) bool {
	redirects := r.GetRedirects()
	r.log.Debug("[ServeHTTP] Redirects: ", fmt.Sprintf("%+v", redirects))

	for _, red := range redirects {
		if red.FromPattern.MatchString(req.URL.Path) {
			r.log.Debug("[ServeHTTP] Redirect from:", req.URL.Path, "to: ", red.To)
			rw.Header().Set("X-Redirected-With", "trf-redirector")
			if r.config.CacheControlMaxAge > 0 {
				rw.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, s-maxage=%d", r.config.CacheControlMaxAge, r.config.CacheControlMaxAge))
			}
			http.Redirect(rw, req, red.To, int(red.Code))
			return true
		}
	}

	return false
}

func (r *TraefikRedirector) ServeLoad(rw http.ResponseWriter, req *http.Request) bool {
	if req.URL.Query().Get(r.config.DebugParameter) == "redirects-load" {
		r.log.Debug("[ServeLoad]")
		res := r.UpdateRedirects()
		if res {
			rw.WriteHeader(http.StatusCreated)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return true
	}

	return false
}

func (r *TraefikRedirector) ServeDebug(rw http.ResponseWriter, req *http.Request) bool {
	if req.URL.Query().Get(r.config.DebugParameter) == "redirects-dump" {
		r.log.Debug("[ServeDebug]")
		jsonResp, err := json.Marshal(r.cache)
		if err != nil {
			r.log.Error("Unable to serialize redirects", err)
			return false
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		_, err = rw.Write(jsonResp)
		if err != nil {
			r.log.Error("Unable to send response", err)
			return false
		}

		return true
	}

	return false
}
