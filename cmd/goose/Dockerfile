FROM golang:1.15.3 AS builder
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/goose cmd/goose
COPY pkg/dbutil pkg/dbutil
ARG SKAFFOLD_GO_GCFLAGS
ENV GCFLAGS=${SKAFFOLD_GO_GCFLAGS}
RUN go build -o goose ./cmd/goose

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/goose /bin/goose
COPY services/todo/migrations/ ./services/todo/migrations/ 
ENV GOTRACEBACK=single
CMD ["/bin/goose"]