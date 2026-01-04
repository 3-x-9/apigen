# apigen

> Generate production-ready CLIs from OpenAPI specifications.

`apigen` is a Go-based tool that converts OpenAPI (Swagger) specs into ergonomic, self-contained command-line interfaces. Instead of copy-pasting `curl` commands or maintaining custom scripts, you get a real CLI with subcommands, flags, config, and autocomplete — automatically.

---

## Why apigen?

APIs are commonly described using OpenAPI, but interacting with them from the terminal is still clunky and repetitive.

`apigen` bridges that gap:

- OpenAPI spec as the single source of truth
- Human-friendly CLI commands generated automatically
- No hand-written clients or boilerplate

---

## Features

- Generate CLIs directly from OpenAPI 3 specifications
- Structured subcommands (e.g. `users list`, `users get`)
- Automatic flag generation from parameters
- Built-in authentication handling
- Config file support
- JSON and table output
- Shell autocomplete (bash, zsh, fish)
- Single static binary with no runtime dependencies

---

## Installation

### Using Go

```bash
go install github.com/yourusername/apigen@latest
