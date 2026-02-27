# DR-002: Validation Rules and Severity

## Status

Accepted

## Context

envlint should not only report missing keys but also validate value formats. Rules need to distinguish between hard errors and informational warnings.

## Decision

| Rule | Detection | Severity |
|------|-----------|----------|
| `missing-key` | Key from example missing in .env | Error |
| `extra-key` | Key in .env but not in example | Warning |
| `required-empty` | Required key has empty value | Error |
| `invalid-url` | Key name contains `URL`, value is not a valid URL | Error |
| `invalid-port` | Key name contains `PORT`, value not 1–65535 | Error |
| `invalid-email` | Key name contains `EMAIL`, no `@` and `.` | Warning |
| `invalid-boolean` | Key name contains `ENABLED`/`ACTIVE`/`IS_`/`DEBUG`, not a bool value | Warning |

Format rules work by key name convention: if the key name contains a specific keyword, the value is validated.

`--strict` promotes all warnings to errors.

## Consequences

- No explicit type annotation needed — convention by key name is sufficient
- URL and port checks are enabled by default, can be disabled via config
- Boolean detection accepts `true/false/1/0/yes/no/on/off`
- Extra keys are only warnings since they are often intentional in practice (local overrides)
