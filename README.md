# Perspective Installer
This repository contains a small helper tool that helps to install Perspective to popular environments like [Docker Compose](https://docs.docker.com/compose/), [Kubernetes](http://kubernetes.io/) and so on.

## Building
To build type:
```
$ govendor sync
$ go build
```
To run tests type:
```
$ go test $(go list ./... | grep -v vendor)
```

## Running
In the simplest form type a list of generators after command name, e.g.:
```
$ perspective-installer docker-compose
```
To output data to custom directory type:
```
$ perspective-installer -dir /path/to/custom/dir docker-compose
```
To see what is going to be done without actually writing data to disk type:
```
$ perspective-installer -dryRun docker-compose
```