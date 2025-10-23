# proxyone

[![goreference](https://pkg.go.dev/badge/github.com/fopina/proxyone.svg)](https://pkg.go.dev/github.com/fopina/proxyone)
[![release](https://img.shields.io/github/v/release/fopina/proxyone)](https://github.com/fopina/proxyone/releases)
[![downloads](https://img.shields.io/github/downloads/fopina/proxyone/total.svg)](https://github.com/fopina/proxyone/releases)
[![ci](https://github.com/fopina/proxyone/actions/workflows/publish-main.yml/badge.svg)](https://github.com/fopina/proxyone/actions/workflows/publish-main.yml)
[![test](https://github.com/fopina/proxyone/actions/workflows/test.yml/badge.svg)](https://github.com/fopina/proxyone/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/fopina/proxyone/graph/badge.svg)](https://codecov.io/github/fopina/proxyone)

A simple HTTP proxy server that supports upstream proxies and no-proxy lists, without breaking TLS.

## Features

- HTTP and HTTPS proxy support
- TLS tunneling (CONNECT method) without breaking encryption
- Configurable upstream proxy
- No-proxy list support
- Configuration via YAML file in standard OS locations

## Configuration

Create a `proxyone.yaml` file in one of these locations:
- `~/.config/proxyone/proxyone.yaml`
- `~/.proxyone/proxyone.yaml`
- `./proxyone.yaml`

Example configuration:

```yaml
upstream_proxy: "http://proxy.company.com:8080"
no_proxy:
  - "localhost"
  - "127.0.0.1"
  - ".local"
  - ".internal"
listen_addr: "localhost:8080"
```

## Usage

```sh
➜  proxyone -h
A simple HTTP proxy server

Usage:
  proxyone [flags]
  proxyone [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Display version

Flags:
  -h, --help   help for proxyone
```

Start the proxy server:

```sh
➜  proxyone
Starting proxy server on localhost:8080
```

## Build

Check out [CONTRIBUTING.md](CONTRIBUTING.md)

### Makefile Targets
```sh
➜  make
bootstrap                      install build deps
build                          build golang binary
clean                          clean up environment
help                           list makefile targets
install                        install golang binary
race                           display test coverage with race
run                            run the app
snapshot                       goreleaser snapshot
test                           display test coverage
```
