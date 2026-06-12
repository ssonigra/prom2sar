# Build stage
FROM golang:1.21 AS builder

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.mod
COPY go.sum go.sum*

# Cache deps before building
RUN go mod download || true

# Copy source code
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o prometheus-dump-operator cmd/main.go

# Final stage
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

WORKDIR /

COPY --from=builder /workspace/prometheus-dump-operator .

USER 65532:65532

ENTRYPOINT ["/prometheus-dump-operator"]
