# Reqium

Reqium is a terminal-first API client built in Go. It starts fast like a lightweight `curl` replacement, but also gives you local workspaces for request history, environments, collections, collection runs, and an interactive terminal UI.

## Features

- Direct HTTP requests from the terminal
- HTTP methods: `GET`, `POST`, `PUT`, `PATCH`, and `DELETE`
- Custom repeatable headers with `--header` or `-H`
- Raw JSON request bodies with `--body` or `-b`
- Request bodies loaded from files with `--body-file` or `-f`
- Response status, headers, body, and duration output
- JSON response pretty-printing by default
- Local request history with replay
- Local environments with `{{variable}}` resolution
- Collections with saved requests
- Collection runner with per-request results
- Interactive terminal UI with history, collections, and environments
- Clean architecture with models, enums, interfaces, app services, and implementations

## Installation

Build from source:

```bash
go build -o bin/reqium ./cmd/reqium
```

Optionally move the generated binary into a directory on your `PATH`.

## Quick Usage

```bash
reqium get https://api.example.com/users
reqium post https://api.example.com/users --header "Content-Type: application/json" --body '{"name":"John"}'
reqium put https://api.example.com/users/1 --body-file ./payload.json
reqium patch https://api.example.com/users/1 -H "Content-Type: application/json" -b '{"name":"Jane"}'
reqium delete https://api.example.com/users/1
```

Run without arguments to open the terminal UI:

```bash
reqium
```

Inside the UI:

```text
n       Compose and send a request
v       Add or update an environment variable
1/2/3   Switch between history, collections, and environments
r       Refresh workspace data
q       Quit
```

In the request composer:

```text
tab                Move to the next field, or complete the highlighted variable when suggestions are open
shift+tab          Move to the previous field
enter              Complete the highlighted variable when suggestions are open
ctrl+left/right    Cycle HTTP method
ctrl+space         Complete variable in terminals that support it
ctrl+@             Complete variable in terminals that emit Ctrl+Space as Ctrl+@
ctrl+s             Send request
esc                Return to dashboard
```

In the environment form:

```text
tab          Move to the next field
shift+tab    Move to the previous field
enter        Save variable and make the environment active
ctrl+s       Save variable and make the environment active
esc          Return to dashboard
```

## Shared Request Flags

```text
-H, --header      Custom header in "Key: Value" format. Repeatable.
-b, --body        Raw request body.
-f, --body-file   Load request body from file.
-t, --timeout     Request timeout in seconds. Default: 30.
    --pretty      Pretty-print JSON responses. Default: true.
    --env         Environment to resolve {{variables}}.
```

## Environments

Create environments and use variables in URLs, headers, and bodies:

```bash
reqium env create dev
reqium env set dev base_url https://api.example.com
reqium env set dev token abc123
reqium env use dev
reqium env list
```

Use variables in requests:

```bash
reqium get "{{base_url}}/users" -H "Authorization: Bearer {{token}}"
reqium post "{{base_url}}/users" -H "Content-Type: application/json" -b '{"name":"John"}'
```

Use a specific environment:

```bash
reqium get "{{base_url}}/users" --env dev
```

## History

Reqium stores executed requests locally.

```bash
reqium history list
reqium history list --limit 50
reqium history show <id>
reqium history replay <id>
```

## Collections

Create collections and save reusable requests:

```bash
reqium collection create users
reqium collection add users list-users GET "{{base_url}}/users"
reqium collection add users create-user POST "{{base_url}}/users" -H "Content-Type: application/json" -b '{"name":"John"}'
reqium collection list
reqium collection show users
```

Run a collection:

```bash
reqium run users --env dev
```

## Local Storage

Reqium stores local data in your user config directory:

```text
<user-config-dir>/reqium/store.json
```

The app layer depends on repository interfaces, so the storage implementation can be replaced later without changing CLI or TUI use cases.

## Project Structure

```text
cmd/reqium/                 Application entrypoint
internal/app/               Use case orchestration
internal/models/            Request, response, history, environment, collection, runner models
internal/enums/             Typed constants such as HTTP methods and runner status
internal/errors/            Shared application errors
internal/interfaces/        Ports for HTTP, formatting, storage, files, and variables
internal/implementations/   HTTP, formatting, filesystem, storage, and variable adapters
internal/cli/               Thin Cobra command layer
internal/tui/               Bubble Tea terminal UI
pkg/version/                Public version package
```

## Developer Commands

```bash
make run
make test
make build
make fmt
```

## Roadmap

- SQLite storage adapter
- Request tabs in the interactive UI
- Full request editor inside the TUI
- Import/export collections
- Postman collection import
- OpenAPI import
- Concurrent request runner
- Response assertions for API testing
- Collection-level pre-request scripts

## License

This project is provided under the MIT License.
