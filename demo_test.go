package plugindemo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/traefik/plugindemo"
)

func TestDemo(t *testing.T) {
	cfg := plugindemo.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := plugindemo.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
        req.Header.Set("Content-Security-Policy-Report-Only", "script-src 'strict-dynamic' 'nonce-fooBar='")
        req.Header.Set("Content-Security-Policy", "script-src 'strict-dynamic' 'nonce-DhcnhD3khTMePgXw'")
        req.Header.Set("irrelevant", "DhcnhD3khTMePgXw")

	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "content-security-policy-report-only", "script-src 'strict-dynamic' 'nonce-fooBar='")
	assertHeader(t, req, "content-security-policy", "script-src 'strict-dynamic' 'nonce-somebodyoncetoldme")
	assertHeader(t, req, "irrelevant", "DhcnhD3khTMePgXw")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
