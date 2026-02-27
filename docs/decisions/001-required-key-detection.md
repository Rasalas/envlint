# DR-001: Required Key Detection

## Status

Accepted

## Context

envlint needs to distinguish whether a missing or empty key is an error or just a warning. A mechanism is needed to mark keys as "required".

## Decision

Three mechanisms, in this priority:

1. **Inline comment `# required`** in `.env.example` — key is required
2. **Non-empty default value** in `.env.example` — key must exist and must not be empty
3. **Config list** in `.envlint.toml` under `[rules.required].keys`

A key without any of these markers and with an empty value in `.env.example` is treated as optional.

## Consequences

- Simple for developers: `# required` annotation is self-explanatory
- Existing `.env.example` files with default values work out of the box
- Config override for teams that don't want to annotate the example file
- Variable references (`$VAR`/`${VAR}`) count as non-empty — a key with a ref value does not trigger a `required-empty` error
