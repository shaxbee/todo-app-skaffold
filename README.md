# Skaffold Todo App

Todo App built in Go with [Skaffold](https://github.com/GoogleContainerTools/skaffold) deployment.

## Features

- [sqlc](https://github.com/kyleconroy/sqlc) - generate database interface and models from migrations and queries
- [openapi-generator](https://github.com/OpenAPITools/openapi-generator) - generate models and clients from OpenAPI spec
- [dockertest](https://github.com/ory/dockertest) - spin up postgres in docker for integration tests

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/).
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Development

### Bootstrap

Bootstrap Kubernetes cluster in docker:

```sh
make bootstrap
```

This will setup local cluster running in Docker using [kind](https://github.com/kubernetes-sigs/kind).

Bootstrap might take some time to deploy Postgres.

### Run

Run the stack in the current terminal:

```sh
bin/skaffold run --force --port-forward
```

Terminating the run will destroy all resources except Postgres.

Skaffold will build Docker images and deploy the stack to the cluster.

### Develop

Continously develop stack:

```sh
bin/skaffold dev --force --port-forward
```

If files change the affected artifacts will be built and re-deployed automatically.

### Test

Following ports are exposed:

- `:8080` API
- `:9000` API Documentation
- `:5432` Postgres

Create a note:

```sh
curl -s -X POST -H 'Content-Type: application/json' http://:8080/api/v1/todo -d '{"title": "Foo", "content": "Bar"}' | jq
```

Get a note:

```sh
curl -s http://:8080/api/v1/todo/:id | jq
```

List notes:

```sh
curl -s http://:8080/api/v1/todo | jq
```

### Cleanup

Destroy the cluster:

```sh
make clean-kind
```

All data in Postgres will be lost.

Full cleanup including downloaded tools:

```sh
make clean
```
