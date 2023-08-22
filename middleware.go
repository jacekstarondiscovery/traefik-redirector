package traefik_redirector

import (
	"context"
	"encoding/json"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/cache"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/client"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/model"
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/provider"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	MaxAge     int64  `json:"maxAge"`
	Endpoint   string `json:"endpoint"`
	Method     string `json:"method"`
	Data       string `json:"data,omitempty"`
	ClientType string `json:"clientType"`
}

func CreateConfig() *Config {
	return &Config{
		MaxAge:     60,
		ClientType: "mock",
		Method:     "GET",
		Endpoint:   "",
		Data:       "",
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

func (r *TraefikRedirector) LoadRedirects() []model.Redirect {
	if r.cache.IsFresh() {
		return r.cache.Redirects
	}

	go func() {
		r.log.Println("[LoadRedirects] Update redirects")
		if !r.cache.Lock() {
			r.log.Println("[LoadRedirects] Cache is locked")
			return
		}

		defer r.cache.Unlock()

		redirects, err := r.provider.GetRedirects(r.config.Method, r.config.Endpoint, r.config.Data)
		if err != nil {
			r.log.Println("[LoadRedirects] Fetch error", err)
			return
		}

		r.log.Println("[LoadRedirects] Update Cache with: ", redirects)
		r.cache.Update(redirects)
	}()

	return r.cache.Redirects
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	logger := log.New(os.Stdout, "[TraefikRedirector]", 0)

	var cl client.HTTPClient
	if config.ClientType == "mock" {
		cl = &client.MockClient{}
	} else {
		cl = &http.Client{}
	}

	return &TraefikRedirector{
		log:      logger,
		config:   config,
		next:     next,
		name:     name,
		cache:    cache.New(config.MaxAge, time.Now),
		provider: provider.New(logger, cl),
	}, nil
}

func (r *TraefikRedirector) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.log.Println("[ServeHTTP] Incoming: ", req.URL.Path)
	if r.ServeDebug(rw, req) {
		return
	}

	if r.ServeRedirect(rw, req) {
		return
	}

	r.next.ServeHTTP(rw, req)
}

func (r *TraefikRedirector) ServeRedirect(rw http.ResponseWriter, req *http.Request) bool {
	redirects := r.LoadRedirects()
	r.log.Println("[ServeHTTP] Redirects: ", redirects)

	for _, red := range redirects {
		if red.FromPattern.MatchString(req.URL.Path) {
			r.log.Println("[ServeHTTP] Redirect from:", req.URL.Path, "to: ", red.To)
			rw.Header().Set("X-Redirected-With", "trf-redirector")
			http.Redirect(rw, req, red.To, int(red.Code))
			return true
		}
	}

	return false
}

func (r *TraefikRedirector) ServeDebug(rw http.ResponseWriter, req *http.Request) bool {
	if req.URL.Query().Get("unicorn") == "redirects" {
		jsonResp, err := json.Marshal(r.cache)
		if err != nil {
			r.log.Println("Unable to serialize redirects", err)
			return false
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")
		_, err = rw.Write(jsonResp)
		if err != nil {
			r.log.Println("Unable to send response", err)
			return false
		}

		return true
	}

	return false
}
