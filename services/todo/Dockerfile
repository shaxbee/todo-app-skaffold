# syntax = docker/dockerfile:1-experimental

FROM golang:1.15.3 AS builder
WORKDIR /opt/build
ARG SKAFFOLD_GO_GCFLAGS
COPY go.* .
COPY pkg pkg
COPY services/todo services/todo
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build ${SKAFFOLD_GO_GCFLAGS} -o todo-service ./services/todo

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/todo-service /bin/todo-service
ENV GOTRACEBACK=single
CMD ["/bin/todo-service"]