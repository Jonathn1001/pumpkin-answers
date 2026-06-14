// Package migrations embeds the golang-migrate SQL files so a single binary can
// run them at startup (and tests can run them against an ephemeral Postgres).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
