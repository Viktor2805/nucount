ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

ENV CGO_ENABLED=0 GOFLAGS=-mod=readonly

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/app ./cmd/main.go

FROM alpine:${ALPINE_VERSION} AS runner
WORKDIR /app

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 appuser

COPY --from=builder /out/app /app/app

USER 10001
EXPOSE 3000
ENTRYPOINT ["/app/app"]
