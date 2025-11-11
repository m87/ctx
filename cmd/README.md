WIP

# Command Blueprint

This document defines the structure and design conventions for creating new commands in the CLI.

## Root Command Structure

All commands related to contexts must be registered under the rootCmd.
This ensures consistent access to all context-related actions from the main command interface.

## Basic Command Pattern

Commands that operate on contexts should follow the structure:

```
ctx action [ "context description" | --ctx-id CONTEXT_ID ]
```

- `"context description"` — a human-readable identifier for the context.
- `--ctx-id CONTEXT_ID` — an explicit context identifier (used when the description is not provided).

Example
```
ctx switch "dev environment"
```

or equivalently:

```
ctx switch --ctx-id 12345
```

## Complex Command Chains

For multi-step or hierarchical operations on contexts, commands should form a chain of subcommands:
```
ctx action [ "context description" | --ctx-id CONTEXT_ID ] subaction params subsubaction params ...
```

Each subsequent subcommand should perform a narrower or more specific operation on the selected context.

Example
```
ctx edit "dev environment" interval --interval-id 123 --start "12-12-2025 12:22:00"
```

This command:

1. Edits the context "dev environment".
2. Modifies the interval with ID 123, updating its start time.

## Design Guidelines

To keep the CLI consistent across contributors, follow these rules:

### Naming

- Top-level after ctx: use verbs (switch, list, create, edit, delete).
- Subactions: can be nouns describing the resource being modified (interval, variable, policy, target).
- All command names should be lowercase.
- Prefer short flags only when obvious (-i for --interval-id is fine, but avoid clever/unclear abbreviations).

### Context selection

Every command that touches a context must accept exactly one of:

- a string description:
```
ctx edit "dev environment" ...
```

- or an explicit ID:
```
ctx edit --ctx-id 123 ...
```

If both are provided, the command should fail with a clear error (don’t guess).

>Note: Context selection by name should be the default for interactive use (human-readable).
Context selection by ID must still be implemented to allow technical or automated execution, such as in scripts

### Parameters

- Prefer --long-flags over positional params for anything non-obvious.
- Date/time params should be explicitly documented (format, timezone if relevant).
- If the command modifies an existing entity (like interval), it must accept an identifier flag (--interval-id, etc.).
- Exception: The context name is not subject to these ID rules.
 - Contexts are the top-level resource and their main identifier is always the context name.
 - Context IDs are generated internally and are primarily intended for non-interactive or automated scenarios.
 - Using the name directly is more natural and user-friendly for human use:
  ```
  ctx edit "prod environment"
  ```
  while IDs remain available for programmatic use:
  ```
  ctx edit --ctx-id 42
  ```

### Help / usage

- Every new command must have a short description (1 sentence) and, if it’s complex, a long description with an example.

- Examples should follow this pattern:
```
# good
ctx switch "dev environment"

# good
ctx edit "dev environment" interval --interval-id 123 --start "12-12-2025 12:22:00"
```

- If the command has mutually exclusive options (e.g. --ctx-id vs "context description"), state it in the help.

### Errors

- Prefer deterministic errors over silent fallback.
 - ✅ “either context description or --ctx-id must be provided”
 - ✅ “--ctx-id and description are mutually exclusive”
 - ❌ “context not found, so I made a new one”

### Consistency with Go/Cobra
- Register under rootCmd → ctxCmd → subcommands.
- Keep command files small: one file per command (e.g. switch.go, editInterval.go).
- Add the example to Example: so --help shows it.