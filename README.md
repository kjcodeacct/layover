![Layover](./assets/layover.png)

# Table of Contents

- [Overview](#overview)
- [Installing & Usage](#installing&usage)
- [Dependencies](#dependencies)
- [Shoutouts & Attribution](#attribution)
- [License](#license)

# Overview
Layover is a TCP & UDP socket proxy intended to help the following scenarios:

a) Debug and log a whily networked application without getting into the weeds with ncat

b) Relaying a port of a non containerized process into a container based network i.e [traefik ❤️](https://github.com/containous/traefik)

### Developers Note 
Please do **not** store debug logs on a live production process, you are essentially logging all traffic, possibly unencrypted.

Logging is intended for *debug* use.

# Installing & Usage

## Installing
While manual installation is *not* recommened for proxying over containers, it is available.

golang 1.14+ is required
```
$ go get -U github.com/kjcodeacct/layover
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
	* the host layover is proxying from, unless specifying to a different host machine use the default

* LAYOVER_PROXYPORT - required:true
	* the port layover is proxying *FROM*

* LAYOVER_PROTOCOL - default:"tcp"
	* IP protocol used by the specified port

* LAYOVER_SERVEPORT default - "8080"
	* the port layover is proxying *TO*

* LAYOVER_DEBUGMODE default - "0"
	* enabled or disable logging
		* 0 - off
		* 1 - basic logging of IP connecting and warnings
		* 2 - full logging including data (please don't use in production)

* LAYOVER_LOGDIR
	* directory to place logfiles created by enabling the DEBUGMODE

# Dependencies
Docker API 1.40+

Golang version 1.14+

# Attribution
Traefik <https://github.com/containous/traefik>

This handy tool <https://github.com/magicmark/composerize>