## softserve

![Softserve Gopher Mascot](softserve-sm.png)

**NOTE:** Flipping this repo public as I've been using it a lot and a coworker asked me for the code. **It's still at the beginning stages**, but it's in a working state. Use wisely and at your own discretion! It's part of my current no-code instrumentation/tool kick.

## Features

| Feature             | Description                                             |
| ------------------- | ------------------------------------------------------- |
| Static file serving | Serve an entire directory of HTML, JS, CSS, and assets  |
| Live reload         | Auto-refresh browser on file changes (injected JS)      |
| API mocking         | YAML-based mock API routes mapped to methods and paths  |
| SSL support         | Built-in HTTPS support with self-signed cert generation |
| Single binary       | No Node, no runtime — just `go run` or build a binary   |
| Graceful shutdown   | Ctrl+C cleanup built-in                                 |

---

## Notes

Mostely tested and used in Linux, but also should work in Windows.

### Frontend Works Out of the Box

Your frontend doesn’t need any special headers, credentials, or port remapping. It simply generates mock responses for api calls — just as if the backend were alive.

This means:

* You can test **error cases** by returning status codes like 401, 403, or 500.
* You can **simulate latency** by adding it later as an advanced feature.
* You can **fully prototype** frontends without scaffolding a backend.

#### No-Touch = No Code

You do **not** need to:

* Add any custom routes in Go
* Define new handlers or switch statements
* Recompile the server

Everything is driven entirely by the contents of `api.yaml`.

## Installation

Clone and build:

```bash
git clone https://github.com/mwiater/softserve.git
cd softserve
go mod tidy
```

`go run cmd/main.go serve --help`

```
Softserve is a lightweight local static file server tailored for frontend development.
It supports automatic browser reloads on file changes, static API mocking, and optional HTTPS serving.

Usage:
  softserve [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Group commands for listing resources
  serve       serve

Flags:
      --api                 enable API mocking
      --api-prefix string   API prefix (default "/api/")
  -h, --help                help for softserve
      --http-port int       HTTP port (default 8080)
      --https-port int      HTTPS port (default 8443)
      --log-level string    log level (default "info")
      --ssl                 enable HTTPS
      --web-root string     directory to serve (default "examples/basic")

Use "softserve [command] --help" for more information about a command.
```

## Configuration

Use command line flags to configure the server. Example (showing defaults):

```bash
go run cmd/main.go serve \
  --web-root examples/basic \
  --ssl=false \
  --http-port 8080 \
  --https-port 8443 \
  --log-level info \
  --api=true \
  --api-prefix /api/
```

---

## Serving Files (Repository Examples)

#### Basic Example

From the root of this repository, simply run with defaults:

```bash
go run cmd/main.go serve
```

```bash
📂 Web root: examples/basic
🌐 Serving HTTP on http://0.0.0.0:8080
```

Then visit: [http://localhost:8080](http://localhost:8080)

#### API Mock Example (with in-memory SSL):

From the root of this repository, simply run with these flags:

```bash
go run cmd/main.go serve --ssl --api --web-root=examples/api01
```

```bash
📂 Web root: examples/api01
🌐 Serving HTTP on http://0.0.0.0:8443
```

Then visit: [https://localhost:8443](https://localhost:8443)

---

### Building and Releasing with Goreleaser

While the examples in this repository use the format `go run cmd/main.go ...`, in the real world you'll build, install and sun it with the `softserve` command.

A minimal `.goreleaser.yaml` is included for building cross platform binaries.
To create local snapshot artifacts without publishing run:

```bash
goreleaser release --snapshot --clean --skip archive
```

Then run:

`./dist/softserve_linux_amd64_v1/softserve`

Global install docs coming soon...

---

## Live Reload

* Automatically injected into `.html` files
* Watches all files recursively in the web root
* WebSocket client reconnects and triggers `location.reload()` on change

---

## Mock API System: How It Works (No-Touch Design)

The **Mock API** feature in *softserve* is a completely *no-touch* system — it does **not** require you to modify or write any Go code to serve dynamic API responses. All responses are defined declaratively in a single `api.yaml` file.

### Goals:

* Serve fake API responses *without running a backend*
* Require **zero changes** to your frontend or backend source code
* Match only on **HTTP method** + **request path**
* Keep things deterministic and inspectable

---

## Development Notes

* All `.html` responses get live-reload JS injected (non-destructively)
* If `index.html` not found, returns 404 (no fallback routing)
* Ignores symlinks for safety
* Logs file changes and server events

---
