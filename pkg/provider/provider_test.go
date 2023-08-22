package provider

import (
	"github.com/jacekstarondiscovery/traefik-redirector/pkg/client"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"regexp"
	"testing"
)

func TestGetRedirects(t *testing.T) {
	t.Run("should return redirects", func(t *testing.T) {
		logger := log.New(os.Stdout, "[TraefikRedirector]", 0)

		p := New(logger, &client.MockClient{})

		result, err := p.GetRedirects("POST", "https://test.dev", "")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Exactly(t, "^/xxx", result[0].From)
		assert.Exactly(t, "https://tvn24.pl", result[0].To)
		assert.Exactly(t, int64(307), result[0].Code)

		rx, _ := regexp.Compile("^/xxx")
		assert.Exactly(t, rx, result[0].FromPattern)
	})
}
