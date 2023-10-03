package traefik_redirector

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeDebug(t *testing.T) {
	t.Run("should respond with debug", func(t *testing.T) {
		next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

		instance, err := New(context.TODO(), next, CreateConfig(), "redirector")
		if err != nil {
			t.Fatalf("Error creating %v", err)
		}

		warmupRecoreder := httptest.NewRecorder()
		warmupReq := httptest.NewRequest(http.MethodGet, "http://localhost/", nil)

		instance.ServeHTTP(warmupRecoreder, warmupReq)
		time.Sleep(1000 * time.Microsecond)

		debugRecorder := httptest.NewRecorder()
		debugReq := httptest.NewRequest(http.MethodGet, "http://localhost/?debug=redirects-dump", nil)

		instance.ServeHTTP(debugRecorder, debugReq)

		body, err := io.ReadAll(debugRecorder.Body)
		assert.NoError(t, err)
		assert.Equal(t, 200, debugRecorder.Result().StatusCode)
		assert.Equal(t, "application/json", debugRecorder.Header().Get("Content-Type"))

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
		assert.Len(t, data["redirects"], 1)
	})
}

func TestBasic(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		called = true
	})

	instance, err := New(context.TODO(), next, CreateConfig(), "redirector")
	if err != nil {
		t.Fatalf("Error creating %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

	instance.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("Invalid return code")
	}
	if called != true {
		t.Fatalf("next handler was not called")
	}
}
