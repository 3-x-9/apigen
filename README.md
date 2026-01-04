apigen

Generate production-ready CLIs from OpenAPI specifications.

apigen is a Go-based tool that turns OpenAPI (Swagger) specs into ergonomic, self-contained command-line interfaces. Instead of copy-pasting curl commands or maintaining custom scripts, you get a real CLI with subcommands, flags, config, and autocomplete — automatically.

Why apigen?

APIs are usually documented with OpenAPI, but interacting with them from the terminal is still painful.

apigen bridges that gap:

📄 OpenAPI in

🧠 Smart command generation

🖥️ A usable CLI out

No hand-written clients. No boilerplate.

Features

Generate CLIs directly from OpenAPI 3 specs

Structured subcommands (e.g. users list, users get)

Automatic flag generation from parameters

Built-in auth handling (API keys, bearer tokens)

Config file support

JSON and table output

Shell autocomplete (bash/zsh/fish)

Single static binary (no runtime dependencies)

Installation
Using Go
go install github.com/yourusername/apigen@latest

From source
git clone https://github.com/yourusername/apigen
cd apigen
go build -o apigen

Quick Start

Given an OpenAPI spec:

paths:
  /users:
    get:
      summary: List users


Generate a CLI:

apigen generate --spec ./openapi.yaml --out ./petstore-cli


Use the generated CLI:

petstore users list
petstore users get --id 123


Under the hood, the CLI handles request construction, auth, and output formatting for you.

How it Works

Load and validate an OpenAPI spec

Parse paths, operations, parameters, and schemas

Generate a Cobra-based CLI command tree

Wire commands to HTTP requests

Output a distributable Go binary

Command Mapping
OpenAPI	CLI
/users	users
GET	list
Path/query parameters	Flags (--id, --limit)
Request body	Flags or stdin
Responses	JSON or table output
Authentication

Currently supported:

API key (header/query)

Bearer tokens

Planned:

OAuth2 device flow

Profile-based auth configs

Auth configuration is stored locally and reused across commands.

Project Structure
apigen/
├── cmd/            # CLI entry points (Cobra)
├── internal/
│   ├── generator/  # OpenAPI → CLI generation logic
│   ├── auth/       # Authentication handlers
│   └── output/     # Output formatting
├── examples/       # Example OpenAPI specs
└── main.go

Status

🚧 Early development

The current focus is:

Core OpenAPI parsing

Command generation

API key auth

JSON output

Expect breaking changes until v1.0.

Roadmap

Pagination helpers

OAuth2 support

Interactive auth setup

Plugin system

Better error messages

More output formats

Why this exists

There are tools that generate API clients and tools that make generic HTTP requests — but few that generate human-friendly CLIs directly from API specs.

apigen focuses on:

Developer experience

Sensible defaults

Extensibility without complexity

Contributing

Contributions are welcome!

Open an issue for bugs or feature requests

Submit a PR with a clear description

Keep changes small and focused

License

MIT
