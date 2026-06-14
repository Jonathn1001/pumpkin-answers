# syntax=docker/dockerfile:1
#
# Single-binary build for the claims-config platform (Plan 4):
#   1. build the React SPA          (node)
#   2. embed it into the Go binary  (golang, //go:embed all:dist)
#   3. ship just the static binary  (distroless, non-root)
#
# Build context is this directory (engineer-challenge). On Railway, set the
# service Root Directory to "engineer-challenge" so these paths resolve.

# ---- Stage 1: build the React SPA -------------------------------------------
FROM node:24-alpine AS web
WORKDIR /web
# Install deps first so this layer is cached until the lockfile changes.
COPY apps/web/package.json apps/web/package-lock.json ./
RUN npm ci
COPY apps/web/ ./
RUN npm run build
# -> /web/dist

# ---- Stage 2: build the Go binary with the SPA embedded ---------------------
FROM golang:1.25-alpine AS build
WORKDIR /src
# Download modules in their own layer for faster incremental builds.
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Drop the built UI into the embed path so `//go:embed all:dist` includes it.
COPY --from=web /web/dist ./cmd/server/dist
# CGO off -> fully static binary that runs on distroless/static.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- Stage 3: minimal runtime ------------------------------------------------
# distroless/static: no shell or package manager, CA certs included, runs as the
# built-in non-root user (uid 65532) — minimal attack surface.
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
# Railway provides PORT at runtime; the server reads it (defaulting to 8080).
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
