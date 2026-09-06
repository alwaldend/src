package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }

type oidcTransport func(*http.Request) (*http.Response, error)

func (f oidcTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestOidcResponseBodiesClosedAndErrorsRedacted(t *testing.T) {
	info := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":{"client_secret":"synthetic-client-secret"}}`)
	}))
	defer info.Close()
	for _, failToken := range []bool{false, true} {
		authBody := &trackedBody{Reader: strings.NewReader(`{"code":"synthetic-code"}`)}
		tokenBody := &trackedBody{Reader: strings.NewReader(`{"id_token":"synthetic-id-token"}`)}
		httpClient := &http.Client{Transport: oidcTransport(func(r *http.Request) (*http.Response, error) {
			status := http.StatusOK
			var body io.ReadCloser = authBody
			if strings.HasSuffix(r.URL.Path, "/token") {
				body = tokenBody
				if failToken {
					status = http.StatusBadRequest
				}
			}
			return &http.Response{StatusCode: status, Body: body, Header: make(http.Header)}, nil
		})}
		config := fixtureVaultConfig(info.URL)
		fetcher := NewOidcFetcher(al.NewVault(config))
		fetcher.httpClient = httpClient
		_, err := fetcher.Get(context.Background(), &injector_proto.Resource{Res: &injector_proto.Resource_Oidc{Oidc: &injector_proto.Oidc{Name: "fixture"}}}, nil)
		if failToken {
			if err == nil {
				t.Fatal("expected token error")
			}
			if strings.Contains(err.Error(), "synthetic-id-token") {
				t.Fatal("response body disclosed")
			}
		} else if err != nil {
			t.Fatal(err)
		}
		if !authBody.closed || !tokenBody.closed {
			t.Fatal("response bodies not closed")
		}
	}
}

func TestOidcRejectsRedirect(t *testing.T) {
	calls := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++ }))
	defer target.Close()
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer source.Close()
	fetcher := NewOidcFetcher(al.NewVault(fixtureVaultConfig(source.URL)))
	_, err := fetcher.Get(context.Background(), &injector_proto.Resource{Res: &injector_proto.Resource_Oidc{Oidc: &injector_proto.Oidc{Name: "fixture"}}}, nil)
	if err == nil {
		t.Fatal("redirect accepted")
	}
	if calls != 0 {
		t.Fatal("redirect destination received request")
	}
}

func TestOidcCancellationStopsRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Drain the POST body so net/http can observe the client's disconnect.
		if _, err := io.Copy(io.Discard, r.Body); err != nil {
			return
		}
		close(started)
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer service.Close()
	// A failing cancellation assertion must not strand httptest.Server.Close.
	defer close(release)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fetcher := NewOidcFetcher(al.NewVault(fixtureVaultConfig(service.URL)))
	done := make(chan error, 1)
	go func() {
		_, err := fetcher.Get(ctx, &injector_proto.Resource{Res: &injector_proto.Resource_Oidc{Oidc: &injector_proto.Oidc{Name: "fixture"}}}, nil)
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("request ignored cancellation")
	}
}

func TestVaultOperationErrorRedactsResponse(t *testing.T) {
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"errors":["synthetic-sensitive-response"]}`)
	}))
	defer service.Close()
	fetcher := NewOpFetcher(al.NewVault(fixtureVaultConfig(service.URL)))
	_, err := fetcher.Get(context.Background(), &injector_proto.Resource{Res: &injector_proto.Resource_Op{Op: &injector_proto.Op{Path: "fixture"}}}, nil)
	if err == nil {
		t.Fatal("expected Vault error")
	}
	if strings.Contains(err.Error(), "synthetic-sensitive-response") {
		t.Fatal("Vault response disclosed in error")
	}
}
