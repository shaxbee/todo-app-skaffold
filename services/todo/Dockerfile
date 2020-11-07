FROM golang:1.15.3 AS builder
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile .
COPY makefiles makefiles
COPY services/todo services/todo
COPY pkg pkg
ARG SKAFFOLD_GO_GCFLAGS
RUN GCFLAGS=${SKAFFOLD_GO_GCFLAGS} make build-todo-service

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/bin/todo-service /bin/todo-service
CMD ["/bin/todo-service"]