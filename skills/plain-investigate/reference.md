# Plain CLI Reference

This file contains reference information for the `plain` CLI tool to help with thread investigation.

## Available Commands

### Search Threads
```bash
plain threads search "query text" [--json]
plain threads search "query" --status=TODO --priority=2
```

### List Threads
```bash
plain threads list [--json]
plain threads list --status=TODO
plain threads list --priority=3 --assignee=user_123
plain threads list --label="Bug"
plain threads list --limit=50 --offset=0
```

### Get Thread Details
```bash
plain threads get <thread-id> [--json]
plain threads get <thread-id> --timeline
```

Accepts either:
- Thread ID: `thread_01h2v3w4x5y6z7a8b9c0`
- Thread URL: `https://app.plain.com/threads/thread_01h2v3w4x5y6z7a8b9c0`

## Common Filters

- **Status**: `TODO`, `DONE`, `SNOOZED`
- **Priority**: `0` (low), `1` (normal), `2` (high), `3` (urgent)
- **Label**: Any label name (case-insensitive)
- **Assignee**: User ID (e.g., `user_123`)

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
        {"id": "label_1", "name": "Bug"},
        {"id": "label_2", "name": "Mobile"}
      ],
      "assignedTo": {"id": "user_456", "name": "Jane Doe"},
      "createdAt": "2026-03-08T10:30:00Z",
      "updatedAt": "2026-03-09T14:20:00Z"
    }
  ],
  "totalCount": 15,
  "hasMore": true
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
    {"id": "label_1", "name": "Bug"},
    {"id": "label_2", "name": "Mobile"}
  ],
  "assignedTo": {
    "id": "user_456",
    "name": "Jane Doe",
    "email": "jane@example.com"
  },
  "createdAt": "2026-03-08T10:30:00Z",
  "updatedAt": "2026-03-09T14:20:00Z",
  "timeline": [
    {
      "type": "status_changed",
      "timestamp": "2026-03-08T10:30:00Z",
      "actor": {"name": "System"},
      "details": {"from": null, "to": "TODO"}
    },
    {
      "type": "assigned",
      "timestamp": "2026-03-08T11:00:00Z",
      "actor": {"name": "John Smith"},
      "details": {"assignee": "Jane Doe"}
    },
    {
      "type": "note_added",
      "timestamp": "2026-03-09T09:00:00Z",
      "actor": {"name": "Jane Doe"},
      "details": {"text": "Investigating iOS version compatibility"}
    }
  ]
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
4. **Check labels for categorization**: Labels often indicate issue type/area
5. **Note assignment patterns**: Certain teams may handle specific issue types
