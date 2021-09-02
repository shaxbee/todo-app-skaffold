FROM golang:1.16.7-buster AS builder
WORKDIR /opt/build
COPY go.* ./
COPY cmd/goose cmd/goose
COPY services/todo/migrations services/todo/migrations
ARG SKAFFOLD_GO_GCFLAGS
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GCFLAGS="${SKAFFOLD_GO_GCFLAGS}" go build -o goose ./cmd/goose

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/goose /bin/goose
COPY services/todo/migrations services/todo/migrations
CMD ["/bin/goose"]