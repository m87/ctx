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

sadsa