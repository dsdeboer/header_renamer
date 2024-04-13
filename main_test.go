package traefik_header_rename

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderRename(t *testing.T) {
	tests := []struct {
		name    string
		rule    Rule
		headers map[string]string
		want    map[string]string
	}{
		{
			name: "no transformation",
			rule: Rule{
				OldHeader: "not-existing",
			},
			headers: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "one transformation",
			rule: Rule{
				OldHeader: "Test",
				NewHeader: "X-Testing",
			},
			headers: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo":       "Bar",
				"X-Testing": "Success",
			},
		},
		{
			name: "Deletion",
			rule: Rule{
				OldHeader: "Test",
			},
			headers: map[string]string{
				"Foo":  "Bar",
				"Test": "Success",
			},
			want: map[string]string{
				"Foo":  "Bar",
				"Test": "",
			},
		},
		{
			name: "one transformation",
			rule: Rule{
				OldHeader: "Test",
				NewHeader: "X-Testing",
			},
			headers: map[string]string{
				"Foo":              "Bar",
				"Test":             "Success",
				"X-Dest-OldHeader": "X-Testing",
			},
			want: map[string]string{
				"Foo":              "Bar",
				"X-Dest-OldHeader": "X-Testing",
				"X-Testing":        "Success",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := CreateConfig()
			cfg.Rules = []Rule{tt.rule}

			ctx := context.Background()
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

			handler, err := New(ctx, next, cfg, "demo-headerrenamerin")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			for hName, hVal := range tt.headers {
				req.Header.Add(hName, hVal)
			}

			handler.ServeHTTP(recorder, req)

			for hName, hVal := range tt.want {
				assertHeader(t, req, hName, hVal)
			}
		})
	}
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
