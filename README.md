# Reqium

Reqium is a terminal API client for sending requests, managing environments, saving history, and running collections without leaving your command line.

Open the interactive interface:

```bash
reqium
```

Or send requests directly:

```bash
reqium get https://api.example.com/users
reqium post https://api.example.com/users -H "Content-Type: application/json" -b '{"name":"John"}'
```

## What You Can Do

- Send `GET`, `POST`, `PUT`, `PATCH`, and `DELETE` requests
- Add custom headers
- Send raw JSON bodies
- Load request body from a file
- Pretty-print JSON responses
- Save local request history
- Create environments with variables like `{{base_url}}`
- Save reusable requests in collections
- Run collections from the terminal
- Use an interactive terminal UI

## Install

From the project folder:

```bash
go install ./cmd/reqium
```

Then run:

```bash
reqium
```

## Interactive UI

Start Reqium:

```bash
reqium
```

Main shortcuts:

```text
n       Create and send a request
v       Add or update an environment variable
1       Show history
2       Show collections
3       Show environments
r       Refresh
q       Quit
```

Request composer:

```text
tab                Next field
shift+tab          Previous field
ctrl+left/right    Change HTTP method
ctrl+s             Send request
esc                Back to dashboard
```

Variable autocomplete:

```text
Type {{ in URL, headers, or body.
Press tab or enter to insert the highlighted variable.
```

## Direct Requests

```bash
reqium get https://api.example.com/users
reqium delete https://api.example.com/users/1
```

With headers:

```bash
reqium get https://api.example.com/users -H "Authorization: Bearer token"
```

With JSON body:

```bash
reqium post https://api.example.com/users \
  -H "Content-Type: application/json" \
  -b '{"name":"John"}'
```

With body from file:

```bash
reqium put https://api.example.com/users/1 -f ./payload.json
```

Common flags:

```text
-H, --header      Header in "Key: Value" format
-b, --body        Raw JSON body
-f, --body-file   Load body from file
-t, --timeout     Timeout in seconds
    --env         Environment to use
    --pretty      Pretty-print JSON responses
```

## Environments

Create an environment:

```bash
reqium env create dev
reqium env set dev base_url https://api.example.com
reqium env set dev token abc123
reqium env use dev
```

Use variables:

```bash
reqium get "{{base_url}}/users" -H "Authorization: Bearer {{token}}"
```

Use a specific environment:

```bash
reqium get "{{base_url}}/users" --env dev
```

## History

```bash
reqium history list
reqium history show <id>
reqium history replay <id>
```

## Collections

Create a collection:

```bash
reqium collection create users
```

Add requests:

```bash
reqium collection add users list-users GET "{{base_url}}/users"
reqium collection add users create-user POST "{{base_url}}/users" \
  -H "Content-Type: application/json" \
  -b '{"name":"John"}'
```

Run a collection:

```bash
reqium run users --env dev
```

## Local Data

Reqium stores history, environments, and collections locally in:

```text
<user-config-dir>/reqium/store.json
```

## License

MIT
