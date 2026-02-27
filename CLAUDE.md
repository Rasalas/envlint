# envlint

CLI tool that validates `.env` files against `.env.example`.

## Tech Stack

- Go 1.25, Cobra CLI, BurntSushi/toml
- No other external dependencies

## Project Structure

- `cmd/` — Cobra commands (root, doctor, init)
- `internal/env/` — .env parser (Entry type, ParseFile)
- `internal/lint/` — Validation rules and linter orchestration
- `internal/config/` — .envlint.toml configuration
- `internal/gitcheck/` — Git/gitignore checks
- `internal/term/` — Color palette and terminal output helpers
- `examples/` — Sample .env + .env.example for testing

## Conventions

- Color palette: Sky Blue (#38BDF8) as primary, differentiated from yeet's Orange
- NO_COLOR support via `init()` in `internal/term/colors.go`
- Output via `term.W` (io.Writer), not stdout directly — enables testability
- Exit codes: 0=pass, 1=errors, 2=config error
- Tests next to code (`_test.go` in same package)
- Patterns follow `../yeet/`
- Documentation always in English

## Development

```bash
go test ./...                              # Tests
go run . --example examples/.env.example --env examples/.env   # Manual test
go run . --format json --example examples/.env.example --env examples/.env
go run . doctor
```
