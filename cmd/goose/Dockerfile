FROM golang:1.15.3 AS builder
WORKDIR /opt/build
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile .
COPY makefiles makefiles
COPY cmd/goose cmd/goose
COPY pkg/dbutil pkg/dbutil
ARG SKAFFOLD_GO_GCFLAGS
RUN GCFLAGS=${SKAFFOLD_GO_GCFLAGS}$ make build-goose

FROM gcr.io/distroless/base-debian10
COPY --from=builder /opt/build/bin/goose /bin/goose
COPY services/todo/migrations/ ./services/todo/migrations/ 
CMD ["/bin/goose"]