package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
)

// uiEngine builds a minimal engine with one real API route plus the SPA mount,
// mirroring how main wires NewRouter + serveUI.
func uiEngine(ui fs.FS) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	mountUI(r, ui)
	return r
}

func TestMountUI(t *testing.T) {
	ui := fstest.MapFS{
		"index.html":        {Data: []byte("<!doctype html><div id=root></div>")},
		"assets/app-123.js": {Data: []byte("console.log('hi')")},
	}
	r := uiEngine(ui)

	cases := []struct {
		name         string
		path         string
		wantStatus   int
		wantBody     string // substring expected in body
		wantCTSubstr string // substring expected in Content-Type (skip if empty)
	}{
		{"root serves index", "/", http.StatusOK, "id=root", "text/html"},
		{"deep client route falls back to index", "/tenants/acme/config", http.StatusOK, "id=root", "text/html"},
		{"real hashed asset served verbatim", "/assets/app-123.js", http.StatusOK, "console.log", "javascript"},
		{"unmatched api stays JSON 404", "/api/missing", http.StatusNotFound, "not_found", "application/json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if w.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body: %s)", w.Code, tc.wantStatus, w.Body.String())
			}
			if !strings.Contains(w.Body.String(), tc.wantBody) {
				t.Fatalf("body = %q, want substring %q", w.Body.String(), tc.wantBody)
			}
			if tc.wantCTSubstr != "" {
				if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, tc.wantCTSubstr) {
					t.Fatalf("Content-Type = %q, want substring %q", ct, tc.wantCTSubstr)
				}
			}
		})
	}
}

// TestSPAFSEmptyIsAPIOnly guards the local-checkout path: with only .gitkeep (no
// index.html) spaFS must report "no UI" so the server stays API-only.
func TestSPAFSEmptyIsAPIOnly(t *testing.T) {
	emptyUI := fstest.MapFS{".gitkeep": {Data: []byte("placeholder")}}
	if _, err := fs.Stat(emptyUI, "index.html"); err == nil {
		t.Fatal("test fixture unexpectedly has index.html")
	}
}
