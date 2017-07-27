package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	if c == nil {
		t.Errorf("Expected client even with nil config")
	}
}

func TestGetTagDigest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method GET; got %q", r.Method)
		}

		accpt := r.Header.Get("Accept")
		if accpt != "application/vnd.docker.distribution.manifest.v2+json" {
			t.Errorf("Expected 'application/vnd.docker.distribution.manifest.v2+json'; got %q, h[0]")
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Etag", "\"sha256:testdigest\"")
	}))

	defer ts.Close()
	c := NewClient(&Config{RootUrl: ts.URL})
	c.GetTagDigest("/test")
}
