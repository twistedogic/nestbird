## Why

NetBird connections drop frequently in flaky environments due to network instability and token refreshes requiring re-authentication. We need an automated watcher service to detect disconnection and reconnect without manual intervention.

## What Changes

- Create `nestbird` — a Go service that monitors NetBird connection status
- Watch NetBird's CLI interface (`netbird status`, `netbird up`) for connection state
- Implement exponential backoff with jitter for reconnection attempts
- Deploy as a systemd service for reliability and automatic restart
- Poll connection status every 5 minutes with max backoff of 5 minutes

## Capabilities

### New Capabilities

- `netbird-watcher`: Core service that monitors NetBird connection state via CLI and automatically reconnects when disconnected

## Impact

- New Go service: `nestbird`
- systemd unit: `nestbird.service`
- Runtime dependency: NetBird CLI (`netbird`)
