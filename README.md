# OwnTracks Recorder
[![](https://img.shields.io/badge/Go-1.19-00ADD8?style=flat&logo=go)](https://golang.org/doc/go1.19)
[![build](https://github.com/hrshadhin/ot-recorder/actions/workflows/build.yml/badge.svg)](https://github.com/hrshadhin/ot-recorder/actions?query=workflow%3ABuild)
[![gosec](https://img.shields.io/github/workflow/status/hrshadhin/ot-recorder/Security?label=%F0%9F%94%91%20gosec&style=flat&color=75C46B)](https://github.com/hrshadhin/ot-recorder/actions?query=workflow%3ASecurity)

[//]: # ([![codecov]&#40;https://codecov.io/gh/hrshadhin/ot-recorder/branch/master/graph/badge.svg?token=N3JVSRO5NZ&#41;]&#40;https://codecov.io/gh/hrshadhin/ot-recorder&#41;)
[![Go Report Card](https://goreportcard.com/badge/github.com/hrshadhin/ot-recorder)](https://goreportcard.com/report/github.com/hrshadhin/ot-recorder)

Store and access data published by OwnTracks apps in (postgres, mysql or sqlite) via REST API

## Architecture
![architecture of ot-recorder](_doc/arch.png)

## System requirements
- OwnTracks app (Android / iOS)
- Postgres/Mysql/Sqlite
- Optional
  - Domain for public access
  - Reverse Proxy(NGINX/HA/Caddy) for TLS, HTTPS
  - Grafana visualization

## Getting started
**WIP**

## Grafana Integration
**WIP**

## Development
- Copy config file `mv _doc/config ./` to root directory and change it
- Local
  ```bash
  make build # build binary
  make version # check binary
  make serve # run the application
  make migrate-up
  make test-unit
  make test-integration # default sqlite
  make test-integration-mysql
  make test-integration-pgsql
  make help # Get all make command
  ```
- Docker
  ```bash
    make docker-build
    make docker-run
    # migrate
    make docker-migrate
  ```
- Visit **`http://localhost:8000`**
- Stop `CTRL + C`

## API's
- Location Ping
- User Last Location

## Docs
- [ERD](_doc/erd.png)
- [API documentation](https://hrshadhin.github.io/projects/ot-recorder/swagger.html)
