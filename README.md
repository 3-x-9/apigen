# apigen

Generate production-ready CLIs from OpenAPI specifications (MVP).

`apigen` is a Go tool that converts OpenAPI 3 specifications into ergonomic, self-contained command-line interfaces. The generated CLIs include subcommands, flags derived from operation parameters, and configuration via `viper`.

> Status: MVP — actively maintained and open to contributions.

## Key ideas

- Generate a CLI command for each OpenAPI path+method
- Auto-generate flags for `path`, `query`, and simple typed parameters
- Use `viper` for config and environment-variable support

## Repository layout (important files)

- `internal/generator/generator.go` — generator core that emits `go.mod`, `config`, `cmd/*`, and `main.go` for the generated CLI
- `cmd/generate.go` — CLI entry that runs the generator
- `main.go` — program entry for `apigen` itself

---

## Quickstart — generate a CLI

1. Build or run `apigen` from the project root:

```bash
# run without installing
go run . generate -s https://petstore3.swagger.io/api/v3/openapi.json -o ./petcli -m petcli

# or build and run
go build -o apigen ./
./apigen generate -s ./path/to/openapi.yaml -o ./outdir -m github.com/yourname/outdir
```

2. Build the generated CLI:

```bash
cd outdir
go build -o mycli ./
# run the generated CLI
./mycli <command> --help
```

Example (petstore):

```bash
go run . generate -s https://petstore3.swagger.io/api/v3/openapi.json -o ./petcli -m petcli
cd petcli
go build -o petcli ./
./petcli get_pet_findByStatus --status available
```

---

## Configuration & environment

- Generated CLIs use `viper` to load configuration. The generator writes a `config.Load(appName string)` that looks for a config file in:
  - `$HOME/.{appName}/config.yaml`
  - `./config.yaml`

- Environment variables are supported. Viper is configured with an env prefix derived from the module/app name (uppercased). For example, with module name `petcli` you can set:

```bash
PETCLI_BASE_URL=https://api.example.com
PETCLI_API_KEY=abcd
```

- Note: `.env` files are not automatically loaded by `viper`. To add `.env` support, include `github.com/joho/godotenv` and call `godotenv.Load()` before `viper` is initialized in `config.Load()`.

---

## Generated CLI behavior

- Flags are generated for operation parameters (`path`, `query`, `header`).
- Path parameters are substituted into the route template.
- Query parameters are URL-encoded using `net/url`/`url.Values`.
- Responses are printed as pretty JSON by default.

Edge cases (arrays, complex serialization, multipart bodies) are not fully implemented in the MVP; contributions welcome.

---

## Contributing

Contributions are welcome. Suggested starters:

- Improve parameter serialization (arrays, explode semantics)
- Add authentication handling (Bearer, API keys) to generated clients
- Add `.env` loading option to generated `config.Load()`
- Add tests for generator output using small OpenAPI fixtures

Steps to contribute:

1. Fork the repo and create a feature branch
2. Run or write tests, update code
3. Open a pull request with a description and rationale

---

## License

This project is MIT licensed. See `LICENSE`.