# MicroLobby - 3rd gen lobbyserver for Warzone 2100

MicroLobby is the next, next gen lobbyserver for Warzone 2100 after [wzlobbserver-ng](https://github.com/Warzone2100/wzlobbyserver-ng).

## Features

- Requires only:
  - Podman/Docker
  - docker-compose
  - task
- Leaves only "~/go" on the OS in podman mode, nothing in docker mode
- Everything in containers
- Automated migrations, migrating on start
- gRPC+Protobuf internal, JSON/XML external
- Argon2-id Hashes
- JWT Tokens
- Integrated RBAC K/V store -> settings/v1
- Loosely coupled Microservices
- Fast to copy&paste a service, easy to start a new one
- Event System as example for IRC/Discord bots
- All communication over NATS. It scales!
- No persistent Data, everything is stored in the DB's

## Basic Architecture

It's written in Golang by using [go-micro.dev/v4](https://go-micro.dev) for simplicity. Transport, Registry and Broker is done over NATS.

The draw.io flowchart for the Architecture:
![Micro Service Architecture](/docs/micro-service-architecture.png)

## Services

### http_proxy

A very simple Proxy to MicroServices. They have to register routes with it over the help of infoservice.

It provides 3 routes, the result will be collected from all microservices:

| METHOD | Route             | AUTH | Description           |
| ------ | ----------------- | ---- | --------------------- |
| GET    | /health           |  n   | Summary health        |
| GET    | /proxy/v1/health  |  y   | Detailed health       |
| GET    | /proxy/v1/routes  |  y   | List of all routes    |

### settings/v1 Service

Basic Key/Value Store with Permissions

### auth/v1 Service

- Give Username + password and you get a JWT back.
- Internaly converts a JWT to a user with roles.

### lobby/v1 Service

Register a game, get list of games and unregister it.

### EMail Service

Sends E-Mails for us.

### OAuth Service

Think it will never be implemented but be part of Profile Service which will be added later.

## Development

### Prerequesits

- [Task](https://taskfile.dev/#/installation)
- docker-compose 1.29+
- podman/docker

### Run

To run this you have to do the following steps:

```bash
git clone https://github.com/pcdummy/microlobby.git
cd microlobby
cp .env.sample .env
task
```

Now enjoy the [health api](http://localhost:8080/health)

### Taskfile

```bash
task -l
```

```text
task: Available tasks for this project:
* build:                Build all containers
* build:protoc:         Generate shared protobuf files and those for a service
* build:tools:          Build tools
* buildup:              Run this when you changed something
* check:toolchain:      Check if you have all tools installed
* default:              Build and run microlobby
* down:                 Stopp all containers
* download:             Download go dependencies
* ps:                   Show running containers
* service:auth_v1:      Build and run/restart auth/v1, restart http_proxy
* service:build:        Build a specific service using env var "SERVICE"
* service:buildrun:     Build and run a service using env var "SERVICE"
* service:http_proxy:   Build and run/restart http_proxy
* service:settings_v1:  Build and run/restart settings/v1, restart http_proxy
* sysctl:               Run sysctl to allow <=80 ports
* up:                   Run all containers
* upgrade:deps:         Update all go dependencies
```

## Authors

- Ren?? Jochum - rene@jochum.dev
- Pastdue (ideas)

## License

Its dual licensed:

- Apache-2.0
- GPL-2.0-or-later
