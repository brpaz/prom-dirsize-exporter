# =========================================
# Build stage
# =========================================
FROM --platform=$BUILDPLATFORM golang:1.22-alpine3.20 as build

ARG TARGETOS
ARG TARGETARCH
ARG BUILD_DATE
ARG GIT_COMMIT
ARG VERSION

ENV GOPROXY=https://proxy.golang.org
ENV GOCACHE=/go/pkg/mod/cache

WORKDIR /app

COPY go.mod go.sum ./

# Set up Go module cache directory and download dependencies
RUN mkdir -p /go/pkg/mod/cache && \
    go mod download && \
    go mod verify

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build \
    -ldflags="-s -w \
    -X /internal/version.BuildDate=${BUILD_DATE} \
    -X github.com/brpaz/prom-dirsize-exporter/internal/version.Version=${VERSION} \
    -X github.com/brpaz/prom-dirsize-exporter/internal/version.GitCommit=${GIT_COMMIT} \
    -extldflags '-static'" -a \
    -o /out/prom-dirsize-exporter ./main.go

# ====================================
# Production stage
# ====================================
FROM alpine:3.20

COPY --from=build /out/prom-dirsize-exporter /bin

RUN chmod +x /bin/prom-dirsize-exporter

ENTRYPOINT ["prom-dirsize-exporter"]
