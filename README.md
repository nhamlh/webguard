# Web interface for Wireguard

A simple, easy to use web interface for Wireguard. It supports SSO authentication (currently Google, Github, Gitlab, Okta are supported) and SCIM2.0 protocol (in development).

## Some screenshots

![Login page](./screenshots/login.png?raw=true "Login page")

![Main page](./screenshots/index.png?raw=true "Main page")

# Installation

Prerequisites:

- Go 1.16
- GCC (required by go-sqlite)

Build Webguard:

```sh
git clone --depth 1 https://github.com/nhamlh/webguard \
    && cd webguard \
    && go build -o ./webguard ./cmd/
```

Run the database schema migration:

```sh
./webguard db migrate
```

Start Webguard:

```sh
./webguard start
```

# Configuration

TBD

# Docker

TBD
