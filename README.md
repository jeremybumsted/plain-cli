# Plain CLI

A command-line interface for working with [Plain](https://plain.com)

## Overview

Plain CLI provides a fast and efficient way to interact with your Plain support workspace without leaving the command line. Manage threads, assignments, labels, and more with simple commands.

Includes a sample agent skill for using with Claude/Codex/Amp/etc.

## Installation

```bash
go install github.com/jeremybumsted/plain-cli@latest
```

Or build from source:

```bash
git clone https://github.com/jeremybumsted/plain-cli
cd plain-cli
go build -o plain
```

## Quick Start

1. Authenticate with Plain:

```bash
plain auth login
```

1. List your threads:

```bash
plain threads list
```

1. Get details on a specific thread:

```bash
plain threads get th_123abc
```

## Features

### Authentication

Manage your Plain API authentication:

- `plain auth login` - Authenticate with Plain API token
- `plain auth status` - Check current authentication status
- `plain auth logout` - Log out and clear credentials

### Thread Management

Core thread operations:

- `plain threads list` - List threads with filters
- `plain threads get <id>` - Get detailed thread information
- `plain threads search <query>` - Search threads by text

### Thread Actions

Perform actions on threads:

- `plain threads assign <thread-id> <user-id>` - Assign thread to a user
- `plain threads unassign <thread-id>` - Unassign a thread
- `plain threads done <thread-id>` - Mark thread as done
- `plain threads todo <thread-id>` - Mark thread as todo
- `plain threads note <thread-id>` - Add a note to a thread
- `plain threads snooze <thread-id>` - Snooze a thread
- `plain threads priority <thread-id>` - Change thread priority

### Label Management

Manage thread labels:

- `plain threads label add <thread-id> <labels...>` - Add labels to a thread
- `plain threads label remove <thread-id> <labels...>` - Remove labels from a thread
- `plain threads label list` - List available label types
- `plain threads label refresh` - Refresh label cache from API

### Field Management

View thread field schemas:

- `plain threads field list` - List available thread fields and their schemas

## Output Formats

All commands support multiple output formats:

- `--json` or `-j` - Output as JSON
- `--quiet` or `-q` - Minimal output (useful for scripting)
- Default: Human-readable table format

Example:

```bash
plain threads list --json
plain threads get th_123abc -q
```

## Configuration

Configure Plain CLI with:

- `--config` flag - Specify custom config file path
- `PLAIN_CONFIG` environment variable - Set default config location

## Thread ID Formats

Thread IDs can be provided in multiple formats:

- Direct ID: `th_123abc`
- Plain URL: The CLI will extract the thread ID automatically

## Examples

### List all todo threads

```bash
plain threads list --status todo
```

### Assign a thread with confirmation

```bash
plain threads assign th_123abc u_456def
```

### Assign a thread without confirmation

```bash
plain threads assign th_123abc u_456def --yes
```

### Add multiple labels to a thread

```bash
plain threads label add th_123abc bug priority-high
```

### Search threads

```bash
plain threads search "login issue"
```

### Get thread details as JSON

```bash
plain threads get th_123abc --json
```

## Version

Check your installed version:

```bash
plain version
```

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## Links

- GitHub: <https://github.com/jeremybumsted/plain-cli>
- Plain: <https://plain.com>
