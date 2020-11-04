FROM golang:1.15.3 AS builder
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/goose cmd/goose
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /bin/goose ./cmd/goose

FROM gcr.io/distroless/base-debian10
COPY --from=builder /bin/goose /bin/goose
COPY services/todo/migrations/ ./services/todo/migrations/ 
CMD ["/bin/goose"]