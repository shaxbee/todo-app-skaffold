# Skaffold Todo App

[![Skaffold Todo App](https://github.com/shaxbee/todo-app-skaffold/workflows/Skaffold%20Todo%20App/badge.svg)](https://github.com/shaxbee/todo-app-skaffold/actions?query=workflow%3A%22Skaffold+Todo+App%22+branch%3A%22master%22)

Todo App built in Go with [Skaffold](https://github.com/GoogleContainerTools/skaffold) deployment.

## Features

- [sqlc](https://github.com/kyleconroy/sqlc) - generate database interface and models from migrations and queries
- [openapi-generator](https://github.com/OpenAPITools/openapi-generator) - generate models and clients from OpenAPI spec
- [dockertest](https://github.com/ory/dockertest) - spin up postgres in docker for integration tests

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/).
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### macOS

Install Docker with brew:

```sh
brew cask install docker
```

Kubectl is included with Docker installation.

## Development

### Bootstrap

Bootstrap Kubernetes cluster in docker:

```sh
make bootstrap
```

This will setup local cluster running in Docker using [kind](https://github.com/kubernetes-sigs/kind).

Bootstrap might take some time to deploy Postgres.

### Run

Run the stack:

```sh
make run
```

Skaffold will build Docker images and deploy the stack to the cluster.
Application will run in the background.

### Debug

Remotely debug the stack:

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
