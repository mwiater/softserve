## 📦 softserve

**Softserve** is a blazing-fast static file server built in Go, designed for front-end developers and full-stack prototypers. It includes:

* ⚡ Static file serving with no config required
* 🔄 Automatic browser reload via WebSockets
* 🧪 Mock API support via YAML definitions
* 🔐 Optional HTTPS with self-signed cert generation

---

## 🚀 Features

| Feature             | Description                                             |
| ------------------- | ------------------------------------------------------- |
| Static file serving | Serve an entire directory of HTML, JS, CSS, and assets  |
| Live reload         | Auto-refresh browser on file changes (injected JS)      |
| API mocking         | YAML-based mock API routes mapped to methods and paths  |
| SSL support         | Built-in HTTPS support with self-signed cert generation |
| Single binary       | No Node, no runtime — just `go run` or build a binary   |
| Graceful shutdown   | Ctrl+C cleanup built-in                                 |

---

### ✅ Frontend Works Out of the Box

Your frontend doesn’t need any special headers, credentials, or port remapping. It simply generates mock responses for api calls — just as if the backend were alive.

This means:

* You can test **error cases** by returning status codes like 401, 403, or 500.
* You can **simulate latency** by adding it later as an advanced feature.
* You can **fully prototype** frontends without scaffolding a backend.

#### 🔄 No-Touch = No Code

✅ You do **not** need to:

* Add any custom routes in Go
* Define new handlers or switch statements
* Recompile the server

Everything is driven entirely by the contents of `api.yaml`.

## 🛠️ Installation

Clone and build:

```bash
git clone https://github.com/YOUR_USERNAME/softserve.git
cd softserve
go build -o bin/softserve cmd/main.go
```

Or just run directly:

```bash
go run cmd/main.go serve
```

---

## 📁 Directory Structure

```
.
├── cmd/
│   └── main.go
│   └── softserve/
│       └── serve.go
│       └── root.go
│       └── list.go
│       └── list_commands.go
├── examples/
│   └── basic/
│   └── api/
├── api.yaml
├── softserve.yaml
├── serve.go
├── watch.go
├── reload.go
├── api.go
├── config.go
├── go.mod
```

---

## 🧾 Configuration

Create a `softserve.yaml` file in your project root:

```yaml
web_root: examples/basic
ssl: false
generate_certs: false
http_port: 8080
https_port: 8443
log_level: info
api: true
api_prefix: /api/
```

---

## 📡 Serving Files

To serve files from `--web-root`:

```bash
go run cmd/main.go serve
```

### Example:

```yaml
web_root: examples/api
```

Then visit: [http://localhost:8080](http://localhost:8080)

---

## 🔄 Live Reload

* Automatically injected into `.html` files
* Watches all files recursively in the web root
* WebSocket client reconnects and triggers `location.reload()` on change

No browser extensions required.

---

## 🧪 Mock API System: How It Works (No-Touch Design)

The **Mock API** feature in *softserve* is a completely *no-touch* system — it does **not** require you to modify or write any Go code to serve dynamic API responses. All responses are defined declaratively in a single `api.yaml` file.

### ✅ Goals:

* Serve fake API responses *without running a backend*
* Require **zero changes** to your frontend or backend source code
* Match only on **HTTP method** + **request path**
* Keep things deterministic and inspectable

---

### 🧱 How It Works Behind the Scenes

1. **API Interception**
   When `api: true` is set in `softserve.yaml`, softserve checks every incoming request to see if:

   * The path starts with the configured `api_prefix` (default: `/api/`)
   * The method and path (e.g. `GET /api/users`) exist in `api.yaml`

2. **Exact Match Lookup**
   Softserve converts the method + path into a key like:

   ```
   GET /api/users
   ```

   It then looks this up in the `api.yaml` map.

3. **Static Response Handling**
   If found, softserve responds with:

   * Status code (e.g. `200`)
   * Any headers (e.g. `Content-Type`)
   * The body text exactly as defined

4. **No Match = 404**
   If no entry exists for a request, the server returns a 404 like a real backend.

---

### 🧪 Example Flow

#### 1. `softserve.yaml`

```yaml
api: true
api_prefix: /api/
```

#### 2. `api.yaml`

```yaml
GET /api/users:
  status: 200
  headers:
    Content-Type: application/json
  body: |
    [
      { "id": 1, "name": "Alice" },
      { "id": 2, "name": "Bob" }
    ]

POST /api/login:
  status: 401
  headers:
    Content-Type: application/json
  body: |
    { "error": "Unauthorized" }
```

#### 3. Run the server:

```bash
go run cmd/main.go serve
```

#### 4. Open your frontend or use curl:

```bash
curl http://localhost:8080/api/users
```

#### ✅ Output:

```json
[
  { "id": 1, "name": "Alice" },
  { "id": 2, "name": "Bob" }
]
```
---

## 🔐 HTTPS (optional)

This tool will autogenerate certs if requested in the `softserve.yaml` config. **If `ssl: true` or `generate_certs: true`, `certs_path` is also required as an absolute path and the `certs_path` must already exist.**

Example `softserve.yaml`:
```
web_root: examples/api
ssl: true                                        # Must be true for SSL
certs_path: /home/matt/projects/softserve/certs  # Required if ssl: true
generate_certs: true                             # Generate self-signed certs
http_port: 8080
https_port: 8443
log_level: info
api: true
api_prefix: /api/
```

Example output when `generate_certs: true`:

```
✅ Config loaded successfully
📂 Web root: examples/api
Checking for existing cert path: '/home/matt/projects/softserve/certs'
  Success: Path is an absolute, existing directory.
🔐 Generated self-signed cert at /home/matt/projects/softserve/certs/
SSL: Loading Cert files:
  >>> /home/matt/projects/softserve/certs/cert.pem
  >>> /home/matt/projects/softserve/certs/key.pem
🔒 Serving HTTPS on https://0.0.0.0:8443

```

---

## 💡 Development Notes

* All `.html` responses get live-reload JS injected (non-destructively)
* If `index.html` not found, returns 404 (no fallback routing)
* Ignores symlinks for safety
* Logs file changes and server events

---

## 🧰 CLI Options (planned)

In addition to `softserve.yaml`, future releases will support:

| Flag               | Description                    |
| ------------------ | ------------------------------ |
| `--web-root`       | Override web\_root from CLI    |
| `--ssl`            | Enable HTTPS                   |
| `--generate-certs` | Generate certs in `certs/`     |
| `--api`            | Enable API mocking             |
| `--api-prefix`     | Override API prefix            |
| `--http-port`      | Port for HTTP (default: 8080)  |
| `--https-port`     | Port for HTTPS (default: 8443) |

---

## ✅ Example Workflow

```bash
cp examples/api/* .
cp api.yaml .
cp softserve.yaml .

go run cmd/main.go serve
```

Visit: [http://localhost:8080](http://localhost:8080)
Call: `curl http://localhost:8080/api/users`

---

## 🧪 Test Cases

* ✅ Auto-refresh when editing HTML
* ✅ Serve nested folders and static assets
* ✅ Exact method/path match for API
* ✅ Missing `index.html` returns 404
* ✅ WebSocket client handles reconnect

---

## 📜 License

TO DO


