# syntax=docker/dockerfile:1.7

############################################
# Builder: Go stable (alpine)
############################################
FROM golang:1.23-alpine AS builder

WORKDIR /src

# Install build tools if needed
RUN apk add --no-cache ca-certificates tzdata upx

# Cache dependencies first
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source
COPY . .

# Build statically linked binary for minimal runtime image
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /out/server ./cmd/server && \
    upx --lzma /out/server || true

############################################
# Runtime: distroless static (nonroot)
############################################
FROM gcr.io/distroless/static:nonroot AS runtime

ENV PORT=8080

WORKDIR /app
COPY --from=builder /out/server /app/server

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/server"]


