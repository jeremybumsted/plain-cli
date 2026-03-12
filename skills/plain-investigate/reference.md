# Plain CLI Reference

This file contains reference information for the `plain` CLI tool to help with thread investigation.

## Authentication & Setup

```bash
plain config          # Interactive setup wizard (first-time or re-configure)
plain config show     # Show current configuration
plain config reset    # Reset configuration
```

API tokens can be generated at: https://app.plain.com/developer/api-keys

## Available Commands

### Search Threads
```bash
plain threads search "query text" [--format=json]
plain threads search "query" --status=TODO --priority=2
plain threads search "query" --label=<label-id>
```

### List Threads
```bash
plain threads list [--format=json]
plain threads list --status=TODO
plain threads list --priority=3 --assignee=user_123
plain threads list --label=<label-id>
plain threads list --limit=50 --offset=0
```

### Get Thread Details
```bash
plain threads get <thread-id> [--format=json]
plain threads get <thread-id> --timeline
```

Accepts either:
- Thread ID: `thread_01h2v3w4x5y6z7a8b9c0`
- Thread URL: `https://app.plain.com/threads/thread_01h2v3w4x5y6z7a8b9c0`

### Modify Thread Status
```bash
plain threads done <thread-id> [--yes]
plain threads todo <thread-id> [--yes]
plain threads snooze <thread-id> [--until=1d] [--yes]
```

Snooze `--until` accepts relative durations (`2h`, `1d`, `3w`) or ISO8601 timestamps.

### Assign Threads
```bash
plain threads assign <thread-id> <user-id> [--yes]
plain threads unassign <thread-id> [--yes]
```

### Thread Notes
```bash
plain threads note <thread-id> [--text="Note text"] [--yes]
# Omit --text to open $EDITOR for longer notes
```

### Thread Priority
```bash
plain threads priority <thread-id> <priority> [--yes]
# Priority accepts: urgent, high, normal, low (or 3, 2, 1, 0)
```

### Thread Labels
```bash
plain threads label list [--refresh] [--archived]    # List available label types
plain threads label add <thread-id> <label...>        # Add labels by name or ID
plain threads label remove <thread-id> <label...>     # Remove labels
plain threads label refresh                           # Refresh label cache from API
```

### Thread Attachments
```bash
plain threads attachments list <thread-id>
plain threads attachments download <attachment-id> [--output=path]
```

### Thread Fields
```bash
plain threads field list    # List available thread field schemas
```

### Articles
```bash
plain articles list [--format=json]
plain articles get <article-id> [--format=json]
```

## Common Filters

- **Status**: `TODO`, `DONE`, `SNOOZED`
- **Priority**: `0`/`low`, `1`/`normal`, `2`/`high`, `3`/`urgent`
- **Label**: Label ID (use `plain threads label list` to find IDs; comma-separated for multiple)
- **Assignee**: User ID (e.g., `user_123`)

## Output Formats

All commands support `--format=<format>` (or global `-j`/`--json` flag):
- `table` (default) — human-readable table
- `json` — structured JSON output for parsing
- `quiet` — minimal output (IDs only)

## Example JSON Output

### Search Results
```json
{
  "threads": [
    {
      "id": "thread_123",
      "title": "Login failure on mobile app",
      "status": "TODO",
      "priority": 2,
      "labels": [
        {"id": "label_1", "labelType": {"name": "Bug"}},
        {"id": "label_2", "labelType": {"name": "Mobile"}}
      ],
      "assignedTo": {"id": "user_456", "fullName": "Jane Doe", "email": "jane@example.com"},
      "createdAt": "2026-03-08T10:30:00Z",
      "updatedAt": "2026-03-09T14:20:00Z"
    }
  ],
  "total": 15
}
```

### Thread Details with Timeline
```json
{
  "id": "thread_123",
  "title": "Login failure on mobile app",
  "description": "User reports unable to login on iOS app version 2.3.0",
  "status": "TODO",
  "priority": 2,
  "labels": [
    {"id": "label_1", "labelType": {"name": "Bug"}},
    {"id": "label_2", "labelType": {"name": "Mobile"}}
  ],
  "assignedTo": {
    "id": "user_456",
    "fullName": "Jane Doe",
    "email": "jane@example.com"
  },
  "createdAt": "2026-03-08T10:30:00Z",
  "updatedAt": "2026-03-09T14:20:00Z",
  "timeline": {
    "entries": [
      {
        "id": "entry_1",
        "timestamp": "2026-03-08T10:30:00Z",
        "actor": {"__typename": "SystemActor"},
        "entry": {
          "__typename": "ThreadStatusTransitionedEntry",
          "previousStatus": null,
          "nextStatus": "TODO"
        }
      },
      {
        "id": "entry_2",
        "timestamp": "2026-03-09T09:00:00Z",
        "actor": {"__typename": "UserActor", "user": {"fullName": "Jane Doe"}},
        "entry": {
          "__typename": "NoteEntry",
          "markdown": "Investigating iOS version compatibility"
        }
      }
    ]
  }
}
```

## Analysis Patterns

When analyzing multiple threads, look for:

1. **Label Clustering**: Which labels appear most frequently?
2. **Status Distribution**: What % are TODO vs DONE vs SNOOZED?
3. **Priority Trends**: Are most issues high priority?
4. **Assignment Patterns**: Who is handling what types of issues?
5. **Time Patterns**: How long do threads typically stay open?
6. **Common Keywords**: What terms appear repeatedly in titles/descriptions?

## Tips for Effective Investigation

1. **Start broad, then narrow**: Begin with a general search, then filter by labels/status
2. **Use timeline for context**: Timeline shows the full story of what happened
3. **Look for related threads**: Extract keywords and search again
4. **Find label IDs first**: Run `plain threads label list` before filtering by label
5. **Note assignment patterns**: Certain teams may handle specific issue types
