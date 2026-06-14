// Command server runs the claims-config platform as a single binary: it mounts
// the /api routes and, when built with the UI (make web / Docker), serves the
// embedded React SPA at / with client-side-route fallback (Plan 4).
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"claimsplatform/internal/configrepo"
	"claimsplatform/internal/httpapi"
	"claimsplatform/internal/seed"
	"claimsplatform/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	setupLogger()
	// Default to release mode so Gin's [GIN-debug] noise stays out of structured
	// logs; set GIN_MODE=debug locally to opt back in.
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fatal("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := configrepo.Migrate(dbURL); err != nil {
		fatal("migrate failed", "err", err)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fatal("db open failed", "err", err) // never log dbURL (contains secret)
	}
	sqlDB, err := db.DB()
	if err != nil {
		fatal("db handle failed", "err", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	repo := configrepo.New(db)
	if err := seed.SeedAll(context.Background(), repo); err != nil {
		fatal("seed failed", "err", err)
	}

	svc := usecase.New(repo)
	router := httpapi.NewRouter(svc)
	serveUI(router) // mounts the embedded SPA when present; API-only otherwise
	slog.Info("api listening", "port", port)
	if err := router.Run(":" + port); err != nil {
		fatal("server failed", "err", err)
	}
}

// setupLogger installs the process-wide slog logger. Defaults to a human-readable
// text handler; set LOG_FORMAT=json for machine-parseable output in production.
func setupLogger() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	var h slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if os.Getenv("LOG_FORMAT") == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(h))
}

// fatal logs an error and exits non-zero (slog has no Fatal; this mirrors log.Fatal).
func fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
