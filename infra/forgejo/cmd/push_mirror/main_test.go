package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExistingHTTPSMirrorIsAdopted(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/repos/alwaldend/other-source/push_mirrors" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte(`[{"remote_name":"remote_mirror_current","remote_address":"https://gitlab.com/alwaldend/other-destination.git","interval":"8h0m0s","sync_on_commit":true}]`))
	}))
	defer server.Close()
	config := mirrorConfig{host: server.URL, repository: "alwaldend/other-source", target: "alwaldend/other-destination", forge: "gitlab", apiToken: "dummy"}
	if err := config.ensure(context.Background(), server.Client()); err != nil {
		t.Fatal(err)
	}
}

func TestNativeRefPreflight(t *testing.T) {
	source := map[string]string{"refs/heads/master": "current", "refs/heads/topic": "new"}
	for _, test := range []struct {
		name string
		dest map[string]string
		want bool
	}{
		{"additive", map[string]string{"refs/heads/master": "current"}, true},
		{"destination-only", map[string]string{"refs/heads/other": "keep"}, false},
		{"divergent", map[string]string{"refs/heads/topic": "keep"}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := verifyRefs(source, test.dest) == nil; got != test.want {
				t.Fatalf("preflight accepted=%t, expected %t", got, test.want)
			}
		})
	}
}
