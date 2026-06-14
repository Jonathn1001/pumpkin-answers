# Build orchestration for the single-binary deploy (Plan 4).
#
#   make web    – build the React SPA and stage it under cmd/server/dist for embedding
#   make build  – build the Go binary with the SPA embedded (bin/server)
#   make run    – build then run locally (needs DATABASE_URL pointing at a Postgres)
#   make test   – run the Go test suite
#   make clean  – remove build artifacts (keeps the dist placeholder)
#
# Local UI development does NOT need `make web`: run `npm run dev` in apps/web and
# the Vite dev server proxies /api to this process (see vite.config.ts).

WEB_DIR   := apps/web
EMBED_DIR := cmd/server/dist
BIN       := bin/server

.PHONY: web build run test clean

web:
	cd $(WEB_DIR) && npm ci && npm run build
	mkdir -p $(EMBED_DIR)
	find $(EMBED_DIR) -mindepth 1 ! -name .gitkeep -delete
	cp -r $(WEB_DIR)/dist/. $(EMBED_DIR)/

build: web
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/server

run: build
	./$(BIN)

test:
	go test ./...

clean:
	rm -rf $(BIN) $(WEB_DIR)/dist
	find $(EMBED_DIR) -mindepth 1 ! -name .gitkeep -delete
