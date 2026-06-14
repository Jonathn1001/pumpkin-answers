// Command server runs the claims-config API for local dev (no SPA embed; that is Plan 4).
package main

import (
	"context"
	"log"
	"os"
	"time"

	"claimsplatform/internal/configrepo"
	"claimsplatform/internal/httpapi"
	"claimsplatform/internal/seed"
	"claimsplatform/internal/usecase"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := configrepo.Migrate(dbURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("db open: %v", err) // never log dbURL (contains secret)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	repo := configrepo.New(db)
	if err := seed.SeedAll(context.Background(), repo); err != nil {
		log.Fatalf("seed: %v", err)
	}

	svc := usecase.New(repo)
	log.Printf("API listening on :%s", port)
	if err := httpapi.NewRouter(svc).Run(":" + port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
