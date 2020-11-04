FROM golang:1.15.3 AS builder
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download
COPY services/todo services/todo
COPY pkg pkg
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /services/todo ./services/todo

FROM gcr.io/distroless/base-debian10
COPY --from=builder /services/todo /services/todo
CMD ["/services/todo"]