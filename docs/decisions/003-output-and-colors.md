# DR-003: Terminal Output and Color System

## Status

Accepted

## Context

envlint needs readable terminal output. The color system should be visually distinct from yeet (orange palette) but follow the same pattern.

## Decision

**Custom color palette with Sky Blue as primary:**

- Primary: `#38BDF8` (Sky Blue)
- Secondary: `#7DD3FC` (light blue)
- Muted: `#78716C` (gray, same as yeet)
- Success: `#4ADE80` (green, same as yeet)
- Danger: `#FF6B6B` (red, same as yeet)
- Warning: `#FBBF24` (yellow, same as yeet)

**NO_COLOR support:** All ANSI codes are set to empty strings in `init()` when the `NO_COLOR` env var is set.

**Output via `term.W`:** All output functions write to `term.W` (default: `os.Stdout`), not directly to stdout. Enables future testability.

**Two output formats:**
- `text` (default): Grouped by category (Missing, Extra, Value Problems) with summary
- `json`: Flat JSON object with `valid`, `total`, `errors`, `warnings`, `issues`

## Consequences

- Visually distinct from yeet
- Same code patterns as yeet, easy to maintain
- CI integration via `--format json` and exit codes
