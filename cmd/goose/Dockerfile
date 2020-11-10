# syntax = docker/dockerfile:1-experimental

FROM golang:1.15.3 AS builder
WORKDIR /opt/build
ARG SKAFFOLD_GO_GCFLAGS
COPY go.* .
COPY pkg/dbutil pkg/dbutil
COPY cmd/goose cmd/goose
COPY services/todo/migrations services/todo/migrations
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build ${SKAFFOLD_GO_GCFLAGS} -o goose ./cmd/goose

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/goose /bin/goose
COPY services/todo/migrations services/todo/migrations
CMD ["/bin/goose"]