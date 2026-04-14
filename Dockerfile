# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# This Dockerfile contains multiple targets.
# Use 'docker build --target=<name> .' to build one.

# ===================================
#
#   Non-release images.
#
# ===================================

# certbuild captures the ca-certificates
FROM docker.mirror.hashicorp.services/alpine:3.22 AS certbuild
RUN apk add --no-cache ca-certificates

# devbuild compiles the binary
# -----------------------------------
FROM golang:1.25.8-alpine@sha256:8e02eb337d9e0ea459e041f1ee5eece41cbb61f1d83e7d883a3e2fb4862063fa AS devbuild
ARG VERSION="dev"
# Set the working directory
WORKDIR /build
RUN go env -w GOMODCACHE=/root/.cache/go-build
# Install dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
COPY . ./
# Build the server
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -ldflags="-s -w -X consul-mcp-server/version.GitCommit=$(shell git rev-parse HEAD) -X consul-mcp-server/version.BuildDate=$(shell git show --no-show-signature -s --format=%cd --date=format:'%Y-%m-%dT%H:%M:%SZ' HEAD)" \
    -o consul-mcp-server ./cmd/consul-mcp-server

# dev runs the binary from devbuild
# -----------------------------------
# Make a stage to run the app
FROM scratch AS dev
ARG VERSION="dev"
# Set the working directory
WORKDIR /server
# Copy the binary from the build stage
COPY --from=devbuild /build/consul-mcp-server .
COPY --from=certbuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Command to run the server (mode determined by environment variables or defaults to stdio)
CMD ["./consul-mcp-server"]

# ===================================
#
#   Release images that uses CI built binaries (CRT generated)
#
# ===================================

# default release image (refereced in .github/workflows/build.yml)
# -----------------------------------
FROM scratch AS release-default
ARG BIN_NAME
# Export BIN_NAME for the CMD below, it can't see ARGs directly.
ENV BIN_NAME=$BIN_NAME
ARG PRODUCT_VERSION
ARG PRODUCT_REVISION
ARG PRODUCT_NAME=$BIN_NAME
# TARGETARCH and TARGETOS are set automatically when --platform is provided.
ARG TARGETOS TARGETARCH
LABEL version=$PRODUCT_VERSION
LABEL revision=$PRODUCT_REVISION
COPY dist/$TARGETOS/$TARGETARCH/$BIN_NAME /bin/consul-mcp-server
COPY --from=certbuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Command to run the server (mode determined by environment variables or defaults to stdio)
CMD ["/bin/consul-mcp-server"]

# ===================================
#
#   Set default target to 'dev'.
#
# ===================================
FROM dev
