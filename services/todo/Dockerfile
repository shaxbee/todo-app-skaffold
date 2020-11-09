FROM golang:1.15.3 AS builder
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download
COPY services/todo services/todo
COPY pkg pkg
ARG SKAFFOLD_GO_GCFLAGS
ENV GCFLAGS=${SKAFFOLD_GO_GCFLAGS}
RUN go build -o todo-service ./services/todo

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/todo-service /bin/todo-service
ENV GOTRACEBACK=single
CMD ["/bin/todo-service"]