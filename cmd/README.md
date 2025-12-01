# Command Design Guidelines for ctx CLI

This document defines the syntax, conventions, and design rules for all commands in the `ctx` command-line interface.
The goal is to maintain a consistent, predictable, and user-friendly CLI API.

 `ctx` manages contexts, which are the primary resource in the system.
All commands follow common CLI design principles inspired by tools such as `git`, `kubectl`, `docker`, and `gh`.


##  General Principles
### 1. Context is the primary resource

Top-level commands operate directly on contexts.
We do not prefix commands with a resource name (e.g., `ctx context list`); instead:

```shell
ctx list
ctx switch Work
ctx delete "Side Project"
```

### 2. Verb-first command structure
```shell
ctx <action> [identifier] [flags] [arguments...]
```
Where identifier may be a name or ID (see rule #7).

### 3. Subcommands for secondary resources
Tags, comments, and intervals are modeled as subcommands:
```shell
ctx tag add Work focus
ctx comment delete Work 7
ctx interval update Work 42 --from ... --to ...
```

### 4. Predictable, POSIX-style flags
Use long flags such as:
```shell
--tag, --comment, --from, --to
```
Short aliases SHOULD be avoided unless they significantly improve usability.

### 5. No “smart” parsing
The grammar MUST remain simple and explicit.
Avoid ambiguous situations or fluent-style chaining.

### 6. Clear, actionable errors
When a command fails, errors should provide:
- the reason,
- what the user can do next,
- usage hints.

### 7. Identifier Rule
Every command that operates on a context or any resource with an ID MUST accept either:
the name of the resource, or
the unique ID of the resource.
Both forms should be interchangeable wherever a resource reference is required.
Examples:

```shell
ctx get Work
ctx get --ctx-id ctx_8f29d12a

ctx comment delete Work --comment-id 7
ctx comment delete --ctx-id ctx_23a7b2f1 --comment-id 7
```

Automations and scripts should rely on IDs, while humans often prefer names.

### 8. Every command MUST provide short and long help
Each command must support two help modes:

Short help 
- One-paragraph explanation of what the command does
- Display of syntax and required arguments

Long help
- Detailed description
- Multiple usage examples
- Edge cases and notes on identifier rules
- Explanation of relevant flags

### 9. Every command MUST support unified output formats
Every command that produces output MUST support the following formats via a shared flag:
- `--output json`
- `--output yaml`
- `--output shell`

#### Default output MUST be human-readable.
This default output should be formatted for readability, using clear labels, indentation, and friendly formatting intended for interactive CLI use.If no format is specified, shell SHOULD be the default for ease of scripting.

Examples:
```shell
ctx get Work                # human-readable output
ctx get Work --output json  # structured output
ctx list --output yaml      # config-style output
ctx tag list Work --output shell   # for scripting
```
Human-readable output SHOULD NOT be used in scripts; json, yaml, or shell SHOULD be selected explicitly for automation.