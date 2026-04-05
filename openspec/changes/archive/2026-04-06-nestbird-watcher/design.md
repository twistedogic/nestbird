## Context

NetBird VPN is used for control plane connectivity. The connection drops periodically due to:
- Flaky network conditions
- Token refreshes that require re-authentication via `netbird up`

The connection is required but transient — when it drops, a simple `netbird up` reconnects it. Manual intervention is tedious and error-prone.

## Goals / Non-Goals

**Goals:**
- Monitor NetBird connection status via CLI (`netbird status`)
- Automatically reconnect when disconnected
- Use exponential backoff with jitter to avoid hammering on transient failures
- Deploy as a systemd service for automatic restart and clean lifecycle management
- Simple, minimal, maintainable code

**Non-Goals:**
- Queue or retry control plane operations (handled separately)
- Metrics or observability beyond basic logging
- Integration with application code
- Support for multiple NetBird peers or complex topology

## Decisions

### CLI vs Library Integration

**Decision:** Use NetBird CLI (`netbird status`, `netbird up`) rather than the client library.

**Rationale:**
- Simpler dependency (single binary vs library)
- No API stability concerns
- Matches how NetBird is already operated

**Alternative:** Use NetBird client library directly for structured responses and better error handling.

### Polling vs Event-Driven

**Decision:** Poll `netbird status` every 5 minutes rather than using events or hooks.

**Rationale:**
- NetBird doesn't expose a stable event mechanism
- 5-minute interval is reasonable for connection monitoring
- Simpler implementation with no external dependencies

**Alternative:** Implement systemd-sleep hooks for suspend/resume detection.

### Backoff Strategy

**Decision:** Exponential backoff with full jitter, 1-minute base, 5-minute cap.

```
Wait time = min(base * 2^attempt, maxBackoff) + random(0, wait/2)
```

| Attempt | Base Wait | With Jitter (±50%) |
|---------|-----------|-------------------|
| 1 | 1m | 30s - 90s |
| 2 | 2m | 1m - 3m |
| 3 | 4m | 2m - 6m (capped at 5m) |
| 4+ | 5m | 2.5m - 5m |

**Rationale:**
- 1m base: aggressive enough to catch brief outages, not excessive
- 5m max: aligns with poll interval, avoids runaway backoff
- Full jitter: distributes load if multiple watchers restart simultaneously

### Go vs Other Language

**Decision:** Implement in Go.

**Rationale:**
- Standard choice for infrastructure tooling
- Easy cross-compilation for Linux targets
- Good ecosystem for CLI tooling and exec management

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         nestbird Watcher                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   main loop:                                                        │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │  for {                                                       │  │
│   │      if connected() {                                        │  │
│   │          sleep(pollInterval)                                 │  │
│   │      } else {                                                │  │
│   │          for !connected() {                                  │  │
│   │              if up() { break }                               │  │
│   │              sleep(backoff())                                │  │
│   │          }                                                   │  │
│   │      }                                                        │  │
│   │  }                                                            │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                                                                     │
│   connected() ──► exec("netbird status") → parse output           │
│   up()        ──► exec("netbird up")                               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
nestbird/
├── main.go           # Entry point, signal handling
├── watcher/
│   ├── watcher.go    # Core loop logic
│   ├── netbird.go    # CLI interaction
│   └── backoff.go    # Exponential backoff with jitter
├── nestbird.service  # systemd unit file
└── Makefile          # Build and install targets
```

## systemd Integration

```ini
[Unit]
Description=Nestbird - NetBird Connection Watcher
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/nestbird
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Notes:**
- `Restart=always` provides a safety net if the watcher crashes
- `RestartSec=5` provides a small delay before restart
- Internal backoff handles retry timing; systemd handles crash recovery

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| `netbird status` output format changes between versions | Log raw output on failure; version pin if needed |
| `netbird up` blocks indefinitely | Set command timeout (e.g., 30s) |
| Network down for extended period | Backoff caps at 5m, watchdog continues |
| systemd restart storm during outage | Jitter prevents synchronized retry attempts |

## Open Questions

- What is the exact output format of `netbird status`? (Need to verify for parsing)
- Does `netbird up` require interactive login, or use cached credentials?
- Should we add systemd-sleep hooks for suspend/resume handling?
