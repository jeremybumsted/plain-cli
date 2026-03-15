# Plain CLI

[![Build status](https://badge.buildkite.com/ef47e410871a642b0a8e25146b3030820aa61b670234180c29.svg)](https://buildkite.com/jeremy-bk/plain-cli)
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

1. Configure Plain CLI (first-time setup):

```bash
plain config
```

This will walk you through:

- API token authentication
- Workspace selection
- Help center selection (for articles)

1. List your threads:

```bash
plain threads list
```

1. Get details on a specific thread:

```bash
plain threads get th_123abc
```

1. Browse help center articles:

```bash
plain articles list
```

## Features

### Configuration

Set up and manage your Plain CLI configuration:

- `plain config` - Interactive setup wizard (configures API token, workspace, and help center)
- `plain config show` - Display current configuration
- `plain config reset` - Clear all configuration

### Authentication

Manage your Plain API authentication:

- `plain auth login` - Authenticate with Plain API token (alternative to `plain config`)
- `plain auth status` - Check current authentication status
- `plain auth logout` - Log out and clear credentials

### Help Center Articles

Access and browse help center documentation (optimized for AI agents):

- `plain articles list` - List all articles in your help center
- `plain articles list --preview` - List articles with content previews (helps identify relevant articles)
- `plain articles get <id>` - Get full article content as clean markdown
- Output formats: markdown (default), JSON, quiet

**AI Agent Integration**: Articles are output in clean markdown format, making them perfect for tools like Claude Code, Cursor, or other AI coding assistants. Agents can search articles with `list --preview` and fetch full content with `get`.

### Thread Management

Core thread operations:

- `plain threads list` - List threads with filters
- `plain threads get <id>` - Get detailed thread information
- `plain threads get <id> --timeline` - Get thread with full timeline history
- `plain threads search <query>` - Search threads by text

#### Date-Based Filtering

Both `plain threads list` and `plain threads search` support flexible date filtering options to help you find threads based on when they were created or last updated.

**Available Date Filter Flags:**

- `--created-after` - Show only threads created after this date
- `--created-before` - Show only threads created before this date
- `--updated-after` - Show only threads updated after this date
- `--updated-before` - Show only threads updated before this date

**Supported Date Formats:**

1. **ISO8601/RFC3339** - Full timestamp with timezone (automatically converted to UTC)
   - Example: `2026-03-15T00:00:00Z`

2. **Date Only** - Simple date format (assumes start of day in UTC)
   - Format: `YYYY-MM-DD`
   - Example: `2026-03-15`

3. **Relative Dates** - Time ago from now
   - `Nd` - N days ago (e.g., `7d` = 7 days ago)
   - `Nw` - N weeks ago (e.g., `2w` = 2 weeks ago)
   - `NM` - N months ago (e.g., `1M` = 1 month ago)
   - `Ny` - N years ago (e.g., `1y` = 1 year ago)

4. **Human-Readable** - Common time references
   - `today` - Start of today
   - `yesterday` - Start of yesterday
   - `last-week` - 7 days ago
   - `last-month` - 1 month ago
   - `last-year` - 1 year ago

All dates are normalized to UTC at the start of the day (00:00:00) for consistency.

### Thread Attachments

View and download attachments from threads:

- `plain threads attachments list <thread-id>` - List all attachments in a thread
- `plain threads attachments download <attachment-id>` - Download a specific attachment
- `plain threads attachments download <attachment-id> --output <path>` - Download with custom filename

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

Configure Plain CLI settings:

### Initial Setup

Run `plain config` for an interactive setup wizard that configures:

- API token (from <https://app.plain.com/developer/api-keys>)
- Workspace selection
- Help center selection

### Config File

Configuration is stored in `~/.config/plain-cli/config.json`

### Environment Variables

Override configuration with environment variables:

- `PLAIN_API_TOKEN` - API token
- `PLAIN_WORKSPACE_ID` - Workspace ID
- `PLAIN_HELP_CENTER_ID` - Help center ID
- `PLAIN_CONFIG` - Custom config file location

### Config File Path

Specify custom config file:

- `--config` flag - Specify config file path for any command
- `PLAIN_CONFIG` environment variable - Set default config location

## Thread ID Formats

Thread IDs can be provided in multiple formats:

- Direct ID: `th_123abc`
- Plain URL: The CLI will extract the thread ID automatically

## Examples

### Configuration

```bash
# First-time setup
plain config

# View current configuration
plain config show

# Reconfigure specific settings
plain config

# Reset all configuration
plain config reset
```

### Help Center Articles

```bash
# List all articles
plain articles list

# List with content previews (helpful for finding relevant articles)
plain articles list --preview

# Get specific article as markdown
plain articles get hca_123abc

# Get article as JSON
plain articles get hca_123abc --format json

# Override help center
plain articles list --help-center-id hc_456def
```

### Thread Management

```bash
# List all todo threads
plain threads list --status todo

# Get thread with full timeline history
plain threads get th_123abc --timeline

# List all attachments in a thread
plain threads attachments list th_123abc

# Download a specific attachment
plain threads attachments download att_789xyz

# Download with custom output path
plain threads attachments download att_789xyz --output /tmp/report.pdf

# Assign a thread with confirmation
plain threads assign th_123abc u_456def

# Assign a thread without confirmation
plain threads assign th_123abc u_456def --yes

# Add multiple labels to a thread
plain threads label add th_123abc bug priority-high

# Search threads
plain threads search "login issue"

# Get thread details as JSON
plain threads get th_123abc --json
```

### Date Filtering Examples

```bash
# Find threads from the last 7 days
plain threads list --created-after 7d

# Find threads from the last week
plain threads list --created-after last-week

# Find threads updated yesterday or later
plain threads list --updated-after yesterday

# Find threads in a specific date range
plain threads list --created-after 2026-03-01 --created-before 2026-03-15

# Find threads created in the last month that are still todo
plain threads list --created-after 1M --status todo

# Search for "bug" in threads from the last 2 weeks
plain threads search "bug" --created-after 2w

# Find high-priority threads updated in the last 3 days
plain threads list --priority urgent --updated-after 3d

# Combine date filters with other filters
plain threads list --created-after last-month --status todo --assignee u_123abc

# Find threads that haven't been updated in over a month
plain threads list --updated-before 1M

# Search for threads in a specific date range with priority filter
plain threads search "payment" --created-after 2026-01-01 --created-before 2026-02-01 --priority high

# Find your assigned threads from the last week
plain threads list --mine --created-after last-week

# Use ISO8601 format for precise timestamps
plain threads list --created-after 2026-03-15T00:00:00Z --created-before 2026-03-15T23:59:59Z
```

## Version

Check your installed version:

```bash
plain version
```

## Releases

New releases are automatically created when a git tag is pushed. The release process:

1. Tag a new version: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. Buildkite automatically builds binaries for all platforms and creates a GitHub release

Download the latest release from [GitHub Releases](https://github.com/jeremybumsted/plain-cli/releases).

Binary archives are provided for:

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64, arm64

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## Links

- GitHub: <https://github.com/jeremybumsted/plain-cli>
- Plain: <https://plain.com>
