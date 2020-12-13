# Skaffold Todo App

[![Skaffold Todo App](https://github.com/shaxbee/todo-app-skaffold/workflows/Skaffold%20Todo%20App/badge.svg)](https://github.com/shaxbee/todo-app-skaffold/actions?query=workflow%3A%22Skaffold+Todo+App%22+branch%3A%22master%22)

Todo App built in Go with [Skaffold](https://github.com/GoogleContainerTools/skaffold) deployment.

## Features

- [sqlc](https://github.com/kyleconroy/sqlc) - generate database interface and models from migrations and queries
- [openapi-generator](https://github.com/OpenAPITools/openapi-generator) - generate models and clients from OpenAPI spec
- [dockertest](https://github.com/ory/dockertest) - spin up postgres in docker for integration tests

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)

### Fetch submodules

```sh
git submodule update --init --recursive
```

### macOS

Install Docker with brew:

```sh
brew cask install docker
```

## Development

### Bootstrap

Setup local cluster running in Docker using [kind](https://github.com/kubernetes-sigs/kind).
Postgres will be installed in this step to avoid recreating database on each run.

```sh
make bootstrap
```

### Run

Skaffold will build Docker images and deploy the stack to the cluster.
Application will run in the background.

```sh
make run
```

### Debug

Run the backend in debugging mode.

```sh
make debug
```

Connect to running backend using todo-service launch configuration in VSCode.

### Develop

Continously develop stack:

```sh
make dev
```

If files change the affected artifacts will be built and re-deployed automatically.

Following ports are exposed when running:

- `:8080` API
- `:9000` API Documentation
- `:5432` Postgres

### Test

Open http://localhost in the browser.
API endpoints can be tested directly from Swagger UI.

```sh
open http://localhost
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
