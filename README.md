# ctx

ctx is a lightweight time-tracking command-line tool and server with an optional Angular web UI. It supports local SQLite mode, remote API mode, and a full web application for managing named contexts and time intervals.

### Interval

- `create interval --context-id <ID> [--start "YYYY-MM-DD HH:MM:SS" | RFC3339] [--end "YYYY-MM-DD HH:MM:SS" | RFC3339] [--status <STATUS>]`
- `list interval [--day YYYY-MM-DD]`
- `edit interval --id <ID> [--context-id <ID>] [--start "YYYY-MM-DD HH:MM:SS" | RFC3339] [--end "YYYY-MM-DD HH:MM:SS" | RFC3339] [--status <STATUS>]`
- `delete interval --id <ID>`
- day summaries (`summary day`)
- context merge (`merge context`)
- local mode (SQLite), remote mode via API (`--remote` or config), and full app mode with Web UI
- output formats: `text`, `json`, `yaml`, `shell`

## Architecture and operation modes

This project can run in 3 modes:

- **CLI local**: `ctx` commands read/write directly to local SQLite
- **CLI remote**: `ctx` commands send REST requests to a server (`--remote` or `remote` in config)
- **Full app (API + Web UI)**: `ctx serve` backend + Angular frontend from `ui/` submodule

## Requirements

- Go 1.25+

## Quick start

### Clone with submodules

```bash
git clone --recurse-submodules <REPO_URL>
cd ctx
```

If the repository is already cloned:

```bash
git submodule update --init --recursive
```

### Build

```bash
go build -o ctx .
```

### Start REST server

```bash
ctx serve --addr :8080
```

### Start full app with Web UI

1. Start backend:

```bash
ctx serve --addr :8080
```

2. Start frontend in a second terminal:

```bash
cd ui
npm install
npm run start
```

3. Open the app:

```text
http://localhost:4200
```

In development mode, UI uses `proxy.conf.json` and forwards `/api` to `http://localhost:8080`.

### Example local usage

```bash
ctx create context --name "Work"
ctx switch --name "Work"
ctx create interval --context-id <CONTEXT_ID>
ctx list context
ctx list interval --day 2026-03-28
ctx summary day --day 2026-03-28
```

### Example remote usage

```bash
ctx --remote http://localhost:8080 create context --name "Work"
ctx --remote http://localhost:8080 list context
ctx --remote http://localhost:8080 switch --name "Work"
ctx --remote http://localhost:8080 free
```

Example for reverse proxy with `/api` prefix:

```bash
ctx --remote http://ctx.example.com/api list context
```

## Configuration

The app reads configuration from `~/.ctx.yaml` and environment variables.

Example `~/.ctx.yaml`:

```yaml
remote: http://localhost:8080
log_level: info
database:
  path: ctx.db
```

Notes:

- `--remote` has higher priority than `remote` in config
- if `--remote` is not provided and `remote` is not configured, commands run locally
- `remote` may include a path prefix, e.g. `https://host/api`

## Global flags

- `--remote, -r` remote server address
- `--output, -o` output format: `text|json|yaml|shell`
- `--config` config file path
- `--verbose, -v` verbose output (include detailed fields and, for `list context`, the intervals). Affects `text`, `json`, `yaml`, and `shell` outputs.

Output examples:

```bash
ctx list context -o json
ctx summary day --day 2026-03-28 -o yaml
ctx list interval --day 2026-03-28 -o shell
```

## Commands

### Context

- `create context --name <NAME>`
- `list context`
- `edit context --id <ID> [--name <NAME>] [--description <TEXT>] [--status <STATUS>]`
- `delete context --id <ID>`
- `switch [--id <ID> | --name <NAZWA>]`
- `free`
- `merge context --source-id <ID> --target-id <ID> [--delete-source=true|false]`

### Interval

- `create interval --context-id <ID> [--start <RFC3339>] [--end <RFC3339>] [--status <STATUS>]`
- `list interval [--day YYYY-MM-DD]`
- `edit interval --id <ID> [--context-id <ID>] [--start <RFC3339>] [--end <RFC3339>] [--status <STATUS>]`
- `delete interval --id <ID>`

### Summary

- `summary day [--day YYYY-MM-DD]`

## REST API (summary)

- `GET /context/`
- `POST /context/`
- `GET /context/{id}`
- `PUT /context/{id}`
- `DELETE /context/{id}`
- `POST /context/switch`
- `POST /context/free`
- `GET /context/{id}/intervals`
- `GET /interval/day/{date}`
- `GET /interval/day/{date}/stats`
- `POST /interval`
- `PUT /interval/{id}`
- `DELETE /interval/{id}`
- `PATCH /interval/{id}/move/{targetId}`

## Docker

Build image:

```bash
docker build -t ctx:latest .
```

Run container:

```bash
docker run --rm -p 8080:8080 -v $(pwd)/data:/data ctx:latest
```

By default, the container starts:

```bash
ctx serve --addr :8080
```

## UI (submodule)

Frontend is added as a Git submodule:

- path: `ui/`
- repo: `https://github.com/m87/ctx-ui`

Useful commands:

```bash
git submodule status
git submodule update --remote --merge ui
```

## License

This project is licensed under `Apache-2.0 license`.
