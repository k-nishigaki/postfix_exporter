[![GitHub release](https://img.shields.io/github/release/hsn723/postfix_exporter.svg?sort=semver&maxAge=60)](https://github.com/hsn723/postfix_exporter/releases)
[![main](https://github.com/Hsn723/postfix_exporter/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/Hsn723/postfix_exporter/actions/workflows/main.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/hsn723/postfix_exporter?tab=overview)](https://pkg.go.dev/github.com/hsn723/postfix_exporter?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/hsn723/postfix_exporter)](https://goreportcard.com/report/github.com/hsn723/postfix_exporter)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hsn723/postfix_exporter)

# Prometheus Postfix exporter

Prometheus metrics exporter for [the Postfix mail server](http://www.postfix.org/).
This exporter provides histogram metrics for the size and age of messages stored in
the mail queue. It extracts these metrics from Postfix by connecting to
a UNIX socket under `/var/spool`. It also counts events by parsing Postfix's
log entries, using regular expression matching. The log entries are retrieved from
the systemd journal, the Docker logs, or from a log file.

## Options

These options can be used when starting the `postfix_exporter`

| Flag                     | Description                                          | Default                           |
|--------------------------|------------------------------------------------------|-----------------------------------|
| `--web.listen-address`   | Address to listen on for web interface and telemetry | `9154`                            |
| `--web.config.file   `   | Path to configuration file that can enable TLS or authentication [(ref)](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md). | `""`                            |
| `--web.telemetry-path`   | Path under which to expose metrics                   | `/metrics`                        |
| `--postfix.showq_path`   | Path at which Postfix places its showq socket        | `/var/spool/postfix/public/showq` |
| `--postfix.logfile_path` | Path where Postfix writes log entries                | `/var/log/mail.log`               |
| `--postfix.logfile_must_exist` | Fail if the log file doesn't exist.            | `true`                            |
| `--postfix.logfile_debug` | Enable debug logging for the log file.              | `false`                           |
| `--log.unsupported`      | Log all unsupported lines                            | `false`                           |
| `--docker.enable`        | Read from the Docker logs instead of a file          | `false`                           |
| `--docker.container.id`  | The container to read Docker logs from               | `postfix`                         |
| `--systemd.enable`       | Read from the systemd journal instead of file        | `false`                           |
| `--systemd.unit`         | Name of the Postfix systemd unit                     | `postfix.service`                 |
| `--systemd.slice`        | Name of the Postfix systemd slice.                   | `""`                              |
| `--systemd.journal_path` | Path to the systemd journal                          | `""`                              |

The `--docker.*` flags are not available for binaries built with the `nodocker` build tag. The `--systemd.*` flags are not available for binaries built with the `nosystemd` build tag.

## Events from Docker

If postfix_exporter is built with docker support, postfix servers running in a [Docker](https://www.docker.com/)
container can be monitored using the `--docker.enable` flag. The
default container ID is `postfix`, but can be customized with the
`--docker.container.id` flag.

The default is to connect to the local Docker, but this can be
customized using [the `DOCKER_HOST` and
similar](https://pkg.go.dev/github.com/docker/docker/client?tab=doc#NewEnvClient)
environment variables.

## Events from log file

The log file is tailed when processed. Rotating the log files while the exporter
is running is OK. The path to the log file is specified with the
`--postfix.logfile_path` flag.

## Events from systemd

Retrieval from the systemd journal is enabled with the `--systemd.enable` flag.
This overrides the log file setting.
It is possible to specify the unit (with `--systemd.unit`) or slice (with `--systemd.slice`).
Additionally, it is possible to read the journal from a directory with the `--systemd.journal_path` flag.

## Build options

By default, the exporter is built without docker and systemd support.

```sh
go build -tags nosystemd,nodocker
```

To build the exporter with support for docker or systemd, remove the relevant build build tag from the build arguments. Note that systemd headers are required for building with systemd. On Debian-based systems, this is typically achieved by installing the `libsystemd-dev` APT package.

```
go build -tags nosystemd
```

## Releases

Signed container images are provided from the GitHub Container Registry (https://github.com/Hsn723/postfix_exporter/pkgs/container/postfix_exporter). The binary included in container images is built without docker and systemd support.

The [Releases](https://github.com/Hsn723/postfix_exporter/releases) page includes signed pre-built binaries for various configurations.

- postfix_exporter binaries are minimal builds (docker and systemd support excluded)
- postfix_exporter_docker binaries have docker support built-in
- postfix_exporter_systemd binaries have systemd support built-in
- postfix_exporter_aio binaries are built with everything included, which can be useful for packaging for systems where the final use-case is not known in advance
