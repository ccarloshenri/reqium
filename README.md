# Reqium

Reqium is a fast, minimal terminal-based API client built in Go. It gives you the essentials of a lightweight Postman-style workflow directly from the command line: send requests, manage headers and payloads, and inspect responses without leaving the terminal.

## Features

- HTTP methods: `GET`, `POST`, `PUT`, `PATCH`, and `DELETE`
- Custom repeatable headers with `--header` or `-H`
- Raw JSON request bodies with `--body` or `-b`
- Request bodies loaded from files with `--body-file` or `-f`
- Response status, headers, body, and duration output
- JSON response pretty-printing by default
- Graceful validation for URLs, headers, body options, JSON payloads, timeouts, and request errors
- Clean layered architecture with dependency injection and testable interfaces

## Installation

Build from source:

```bash
go build -o bin/reqium ./cmd/reqium
```

Optionally move the generated binary into a directory on your `PATH`.

## Usage

```bash
reqium get https://api.example.com/users
reqium post https://api.example.com/users --header "Content-Type: application/json" --body '{"name":"John"}'
reqium put https://api.example.com/users/1 --body-file ./payload.json
reqium patch https://api.example.com/users/1 -H "Content-Type: application/json" -b '{"name":"Jane"}'
reqium delete https://api.example.com/users/1
```

Shared flags:

```text
-H, --header      Custom header in "Key: Value" format. Repeatable.
-b, --body        Raw request body.
-f, --body-file   Load request body from file.
-t, --timeout     Request timeout in seconds. Default: 30.
    --pretty      Pretty-print JSON responses. Default: true.
```

## Project Structure

```text
cmd/reqium/                Application entrypoint
internal/app/              Use case orchestration
internal/domain/           Request and response models, validation, domain errors
internal/interfaces/       Ports for HTTP client, formatter, and file reader
internal/infrastructure/   HTTP, formatting, and filesystem adapters
internal/cli/              Thin Cobra command layer
pkg/version/               Public version package
```

## Developer Commands

```bash
make run
make test
make build
make fmt
```

## Roadmap

- Save request history locally
- Collections
- Environment variables
- Export/import requests
- OpenAPI import
- Concurrent request runner
- Response assertions for API testing

## License

This project is provided under the MIT License.
