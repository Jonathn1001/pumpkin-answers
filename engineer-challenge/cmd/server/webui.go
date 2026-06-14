package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// distFS holds the built SPA (apps/web/dist, copied to ./dist by `make web`)
// embedded into the binary so a single artifact serves both the API and the UI
// in production. In a plain checkout ./dist contains only .gitkeep, so spaFS
// reports "no UI" and the server runs API-only — during local development the
// Vite dev server serves the UI and proxies /api to this process instead.
//
//go:embed all:dist
var distFS embed.FS

// spaFS returns the embedded dist rooted at its top level, plus whether a real
// UI build is present (index.html exists). When false, callers skip UI mounting.
func spaFS() (fs.FS, bool) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, false
	}
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		return nil, false
	}
	return sub, true
}

// serveUI mounts the embedded SPA onto r when a build is present; otherwise it
// leaves the server API-only. Logged either way so the deployment mode is clear.
func serveUI(r *gin.Engine) {
	ui, ok := spaFS()
	if !ok {
		slog.Info("web UI not embedded; serving API only")
		return
	}
	mountUI(r, ui)
	slog.Info("web UI embedded; serving SPA")
}

// mountUI registers a catch-all that serves the SPA from ui. Real files are
// served as-is (http.FileServer sets content-type, ETag, and handles ranges);
// any other non-API path falls back to index.html so client-side routes resolve
// on a hard refresh. Unmatched /api paths still return JSON 404, never the HTML
// shell — otherwise a mistyped API call would look like a 200 to clients.
func mountUI(r *gin.Engine, ui fs.FS) {
	fileServer := http.FileServer(http.FS(ui))
	r.NoRoute(func(c *gin.Context) {
		reqPath := c.Request.URL.Path
		if reqPath == "/api" || strings.HasPrefix(reqPath, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "not_found", "message": "route not found"}})
			return
		}
		clean := strings.TrimPrefix(path.Clean(reqPath), "/")
		if clean == "" {
			clean = "index.html"
		}
		if _, err := fs.Stat(ui, clean); err != nil {
			// Not a real asset → treat as a client-side route and serve the shell.
			c.Request.URL.Path = "/"
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}
