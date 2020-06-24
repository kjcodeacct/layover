![Layover](./assets/layover.png)

# Table of Contents

- [Overview](#overview)
- [Examples](#examples)
- [Dependencies](#dependencies)
- [Shoutouts & Attribution](#attribution)
- [License](#license)

# Overview
Layover is a TCP & UDP socket proxy intended to help the following scenarios:

a) Debug and log a whily networked application without getting into the weeds with ncat

b) Relaying a port of a non containerized process into a container based network i.e [traefik ❤️](https://github.com/containous/traefik)

### Developers Note 
Please do store debug logs on a live or production process, you are essentially logging all traffic, most likely unencrypted.

Logging is intended for debug use.

# Examples
Coming Soon

## Building
Go version 1.14+ required

Manual Compilation
```
$ go build
```

# Installing & Usage

## Manual
While manual installation is *not* recommened for proxying over containers, it is available.

golang 1.14+ is required
```
$ go get -U github.com/kjcodeacct/layover
```

### Usage

#### Minimal
```
$ export LAYOVER_SERVEPORT=8080
$ export LAYOVER_PROXYPORT=8081
$ layover
```

## Docker

While it is recommended to use docker-compose for simplicity, docker run is available.

Docker Run
```
docker run -d --restart-always -p 8080:8080 kjcodeacct/layover -e LAYOVER_PROXYPORT=8081 
```

Docker Compose
```
docker-compose up layover
```

# Dependencies
Docker API 1.40+

Golang version 1.14+

# Attribution
Traefik <https://github.com/containous/traefik>

This handy tool <https://github.com/magicmark/composerize>