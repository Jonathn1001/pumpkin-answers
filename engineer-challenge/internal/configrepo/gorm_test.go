package configrepo_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"claimsplatform/internal/configrepo"
	"claimsplatform/internal/configrepo/repotest"
	"claimsplatform/internal/domain"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sharedDSN string

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctr, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("claims"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Println("skipping configrepo integration tests (docker unavailable):", err)
		os.Exit(0)
	}
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}
	sharedDSN = dsn
	code := m.Run()
	_ = ctr.Terminate(ctx)
	os.Exit(code)
}

func newGormRepo(t *testing.T) domain.ConfigurationRepository {
	t.Helper()
	if sharedDSN == "" {
		t.Skip("no database")
	}
	db, err := gorm.Open(gormpg.Open(sharedDSN), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// Isolate each sub-test: drop & recreate the schema, then re-migrate.
	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		t.Fatalf("reset schema: %v", err)
	}
	if err := configrepo.Migrate(sharedDSN); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return configrepo.New(db)
}

func TestGormRepoSatisfiesContract(t *testing.T) {
	repotest.Contract(t, newGormRepo)
}
