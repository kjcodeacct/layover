![Layover](./assets/layover.png)

---
![License](https://img.shields.io/github/license/kjcodeacct/layover)
[![Go Report Card](https://goreportcard.com/badge/github.com/kjcodeacct/layover)](https://goreportcard.com/report/github.com/kjcodeacct/layover)
[![Build Status](https://cloud.drone.io/api/badges/kjcodeacct/layover/status.svg)](https://cloud.drone.io/kjcodeacct/layover)
[![Docker Build Status](https://img.shields.io/docker/build/kjcodeacct/layover)](https://hub.docker.com/repository/docker/kjcodeacct/layover)


# Table of Contents

- [Table of Contents](#table-of-contents)
- [Overview](#overview)
- [Quick Start](#quick-start)
- [Installing & Usage](#installing--usage)
	- [Installing](#installing)
		- [Binary Releases](#binary-releases)
		- [Go Get](#go-get)
	- [Building](#building)
		- [Binaries](#binaries)
		- [Docker](#docker)
		- [CI/CD](#cicd)
	- [Usage](#usage)
		- [CLI](#cli)
		- [Docker](#docker-1)
		- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Attribution](#attribution)

# Overview
Layover is a TCP & UDP socket proxy intended to help the following scenarios:

* Relaying a port of a non containerized process into a container based network i.e [traefik ❤️](https://github.com/containous/traefik)
* Debug and log a whily networked application

# Quick Start
For a very quick deployment please do the following
* set your proxy port variable 'LAYOVER_PROXYPORT' in the docker-compose.yml
* if you are using traefik, modify the host for 'traefik.http.routers.layover.rule', for more information
go [here](https://docs.traefik.io/user-guides/docker-compose/basic-example/)

* run the following
```
$ docker-compose up -d --force-recreate --build layover
```

# Installing & Usage

## Installing
While installation is *not* recommened for proxying over containers, it is available.

### Binary Releases
Binary releases are available in the github releases page found [here](https://github.com/kjcodeacct/layover/releases)

### Go Get
golang 1.14+ is required
set GO111MODULE=on

```
$ go get -u github.com/kjcodeacct/layover
```

## Building


### Binaries
If you would like to manually build binaries available in the [releases page](https://github.com/kjcodeacct/layover/releases), run the following.
```
$ make binaries
```

### Docker
If you would like to build your own release of layover, please see the Dockerfile for a local build
and run the following for a successful image build.

```
$ make docker
```

### CI/CD
If you want to view steps used by <drone.io> for automated builds please view [.drone.yml](.drone.yml)

## Usage

### CLI
```
$ export LAYOVER_SERVEPORT=8081
$ export LAYOVER_PROXYPORT=8080
$ layover proxy
```

### Docker

Docker Run
```
$ docker run -e LAYOVER_SERVEPORT=8081 -e LAYOVER_SERVEHOST=0.0.0.0 -e LAYOVER_PROXYHOST=172.17.0.1 -e LAYOVER_PROXYPORT=8080 --publish 8081:8081 kjcodeacct/layover
```

Docker Compose
```
$ docker-compose up layover
```

### Configuration
Below is a complete list of configuration for more complex needs

Below are the available flags and their configuration.
Flags can be provided either by the CLI flag or via ENV variables
| Flag      | Env Variable        | Type   | Description                                                                                                                      | Default   |
| :-------- | :------------------ | :----- | :------------------------------------------------------------------------------------------------------------------------------- | :-------- |
| proxyhost | `LAYOVER_PROXYHOST` | string | the host layover is proxying from, unless specified to a different host machine uses the default                                 | 0.0.0.0   |
| proxyport | `LAYOVER_PROXYPORT` | int    | the port layover is proxying *FROM*, this is *typically* the port not in the container system                                    | 8081      |
| serveport | `LAYOVER_SERVEPORT` | int    | the port layover is proxying *TO* and is serving, if running in a container typically does *not* need to be specified            | 8080      |
| servehost | `LAYOVER_SERVEHOST` | string | the host layover is proxying *TO* and is serving, if running in a container typically this needs to be specified to '172.17.0.1' | 127.0.0.1 |

# Dependencies
Docker API 1.40+

Golang version 1.14+

# Attribution
Portions of TCP & UDP network handling was used from:
Traefik <https://github.com/containous/traefik>