# Skaffold Todo App

TODO App built in Go with [Skaffold](https://github.com/GoogleContainerTools/skaffold) deployment.

## Prerequisites

Docker is used for running the project.
Install it using [Get Docker](https://docs.docker.com/get-docker/) guide.

## Development

### Bootstrap

Bootstrap Kubernetes cluster in docker:

```sh
make bootstrap
```

This will download [kind](https://github.com/kubernetes-sigs/kind) and create local cluster running in Docker.

Bootstrap might take some time to deploy Postgres.

### Run

Run the stack in the current terminal:

```sh
bin/skaffold run --force --port-forward
```

Terminating the run will destroy all resources except Postgres.

Skaffold will build Docker images and deploy the stack to the cluster.

Ports:

- `:8080` API
- `:5432` Postgres

### Develop

Continously develop stack:

```sh
bin/skaffold dev --force --port-forward
```

If files change affected artifacts will be built and re-deployed automatically.

Ports remain the same as in [Run](#run).

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
