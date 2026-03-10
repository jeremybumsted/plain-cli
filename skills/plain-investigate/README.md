# Plain Investigation Skill

A semi-automated skill for investigating support threads in Plain.

## Purpose

This skill helps you search, analyze, and understand support issues by:
- Finding relevant threads based on keywords or thread IDs
- Analyzing patterns across multiple threads
- Providing actionable insights and recommendations
- Showing timeline history and cross-thread patterns

## Usage

Invoke the skill with `/plain-investigate` followed by:

### Search by Keywords
```
/plain-investigate login failures
/plain-investigate payment issues with stripe
/plain-investigate mobile app crashes
```

### Investigate Specific Thread
```
/plain-investigate thread_01h2v3w4x5y6z7a8b9c0
/plain-investigate https://app.plain.com/threads/thread_123
```

### General Investigation
```
/plain-investigate
```
(The skill will ask what you want to investigate)

## What It Does

1. **Searches** for relevant threads based on your query
2. **Analyzes** thread details, status, priority, labels, and assignments
3. **Identifies patterns** across multiple related threads
4. **Shows timeline** of key events and status changes
5. **Provides insights** about common themes and trends
6. **Suggests next steps** for further investigation

## Example Session

```
You: /plain-investigate login failures

Skill: I'll search for threads about: login failures
      [Runs: plain threads search "login failures" --json]

      Found 12 threads related to login failures
      Common labels: Bug, Mobile, Authentication
      Status: 5 todo, 4 done, 3 snoozed
      Priority: 3 urgent, 6 high, 3 normal

      [Shows detailed analysis of top threads]

      ## Pattern Analysis
      The "login failure" issue appears to cluster around:
      - Mobile app users (label: "mobile")
      - High priority (8 of 12 threads are priority 2-3)

      ## Suggested Next Steps
      1. Check for similar mobile-related issues:
         `plain threads list --label=mobile --status=TODO`

      Would you like me to investigate a specific thread more deeply?
```

## Features

- **Semi-automated**: Explains what it's doing before executing
- **Structured output**: Clear formatting with tables and lists
- **Pattern recognition**: Identifies trends across threads
- **Follow-up suggestions**: Recommends next investigation steps
- **Timeline analysis**: Shows full event history for threads
- **Error handling**: Gracefully handles missing threads or auth issues

## Requirements

- `plain` CLI tool installed and authenticated (`plain auth login`)
- Valid API token with read access to threads

## Customization

The skill focuses on **read operations only**. It will:
- Search and list threads
- Get thread details with timeline
- Analyze patterns and provide insights

It will NOT automatically:
- Modify thread status
- Assign threads
- Add labels or notes
- Make any write operations

This keeps investigations safe and non-destructive.
