![Layover](./assets/layover.png)

# Table of Contents

- [Overview](#overview)
- [Quick Start](#quickstart)
- [Installing & Usage](#installing&usage)
- [Dependencies](#dependencies)
- [Shoutouts & Attribution](#attribution)
- [License](#license)

# Overview
Layover is a TCP & UDP socket proxy intended to help the following scenarios:

a) Debug and log a whily networked application without getting into the weeds with ncat

b) Relaying a port of a non containerized process into a container based network i.e [traefik ❤️](https://github.com/containous/traefik)

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

### CLI
golang 1.14+ is required
```
$ go get -U github.com/kjcodeacct/layover
```

### Docker
If you would like to build your own release of layover, please see the /src/Dockerfile for a local build
and run the following for a successful image build.

```
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o layover
$ docker build . -t layover:test
$ docker run -e LAYOVER_PROXYPORT=8090 --publish 8080:8080 --name test1 layover:test
```

## Usage

### CLI
```
$ export LAYOVER_SERVEPORT=8080
$ export LAYOVER_PROXYPORT=8081
$ layover
```

### Docker
While it is recommended to use docker-compose for simplicity, docker run is available.

Docker Run
```
docker run -d --restart-always -p 8080:8080 kjcodeacct/layover -e LAYOVER_PROXYPORT=8081 
```

Docker Compose
```
docker-compose up layover
```

### Configuration
Below is a complete list of configuration for more complex needs

* LAYOVER_PROXYHOST - default:"0.0.0.0"
	* the host layover is proxying from, unless specifying to a different host machine uses the default

* LAYOVER_PROXYPORT - required:true
	* the port layover is proxying *FROM*
    * this is *typically* the port not in the container system

* LAYOVER_PROTOCOL - default:"tcp"
	* IP protocol used by the specified port
    * options available 
        * "tcp"
        * "udp"

* LAYOVER_SERVEPORT default - default:"8080"
	* the port layover is proxying *TO* and is serving
    * if running in a container typically does *not* need to be specified

* LAYOVER_DEBUGMODE default - "0"
	* options available
		* 0 - off
		* 1 - basic logging of IP connecting and warnings
		* 2 - full logging including data (please don't use in production)

* LAYOVER_LOGDIR
	* directory to place logfiles created by enabling the LAYOVER_DEBUGMODE

### Developers Note 
Please do **not** store debug logs on a live production process, you are essentially logging all traffic, possibly unencrypted.

Logging is intended for *debug* use.

# Dependencies
Docker API 1.40+

Golang version 1.14+

# Attribution
Traefik <https://github.com/containous/traefik>

This handy tool <https://github.com/magicmark/composerize>