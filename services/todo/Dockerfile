FROM golang:1.17.0-buster AS builder
# copy source files
WORKDIR /opt/build
COPY go.* ./
COPY pkg pkg
COPY services/todo services/todo
# build
ARG SKAFFOLD_GO_GCFLAGS
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    GCFLAGS="${SKAFFOLD_GO_GCFLAGS}" go build -o todo-service ./services/todo

# create runtime container from minimal image containing libc, tzdata and root-ca
FROM gcr.io/distroless/base-debian10
# copy binary built in builder container
COPY --from=builder /opt/build/todo-service /bin/todo-service
# setup GOTRACEBACK, it will be overriden by Skaffold when debugging
ENV GOTRACEBACK=single
CMD ["/bin/todo-service"]