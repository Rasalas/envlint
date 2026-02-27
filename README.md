# envlint

Validate `.env` files against `.env.example`. Catches missing keys, empty required values, invalid formats (URLs, ports, emails, booleans), and extra keys. Pre-commit-hook- and CI-ready with clean terminal output.

## Install

```bash
go install github.com/rasalas/envlint@latest
```

## Usage

```bash
# Lint .env against .env.example (defaults)
envlint

# Custom paths
envlint --example .env.example --env .env.local

# Strict mode (warnings become errors)
envlint --strict

# JSON output for CI
envlint --format json

# Only show errors
envlint --quiet
```

### Subcommands

```bash
# Check project setup
envlint doctor

# Generate .env.example from existing .env
envlint init
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | Errors found |
| 2 | Configuration error (missing files) |

## Validation Rules

| Rule | Trigger | Severity |
|------|---------|----------|
| `missing-key` | Key from example missing in .env | Error |
| `extra-key` | Key in .env but not in example | Warning |
| `required-empty` | Required key has empty value | Error |
| `invalid-url` | Key contains `URL`, value is not a URL | Error |
| `invalid-port` | Key contains `PORT`, value not 1–65535 | Error |
| `invalid-email` | Key contains `EMAIL`, invalid format | Warning |
| `invalid-boolean` | Key contains `ENABLED`/`ACTIVE`/`IS_`, not a bool | Warning |

### Required Keys

A key is required when any of these apply:

1. Inline comment `# required` in `.env.example`
2. Non-empty default value in `.env.example`
3. Listed in `[rules.required]` in `.envlint.toml`

Variable references (`$VAR` / `${VAR}`) count as non-empty.

## Configuration

Optional `.envlint.toml` in the project root:

```toml
example = ".env.example"
envFiles = [".env"]

[rules]
requireAll = true
noExtra = false
strictUrls = true
strictPorts = true

[rules.required]
keys = ["DATABASE_URL", "API_KEY"]

[rules.ignore]
keys = ["OPTIONAL_DEBUG_FLAG"]
```

## Pre-commit Hook

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/rasalas/envlint
    rev: v0.1.0
    hooks:
      - id: envlint
```

## Example Output

```
  envlint · .env vs .env.example

  Missing keys

  ✗ DATABASE_URL  (required)
  ✗ REDIS_URL

  Extra keys

  ! LEGACY_FLAG
  ! OLD_API_KEY

  Value problems

  ✗ API_URL — invalid URL format
  ✗ APP_PORT — must be 1-65535, got "abc"
  ! SMTP_PORT — empty value

  ✓ 12 of 18 keys valid · 4 errors · 2 warnings
```

Supports `NO_COLOR=1` for plain output.
