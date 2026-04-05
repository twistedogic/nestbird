## Why

Some environments require a setup key to authenticate with the NetBird management server on `netbird up`. Currently nestbird always runs `netbird up` with no arguments, making it unusable in environments where re-authentication requires a setup key file.

## What Changes

- `netbird up` optionally passes `--setup-key-file=<path>` when a path is configured
- Setup key file path is configurable via CLI flag (`--setup-key-file`) or environment variable (`NETBIRD_SETUP_KEY_FILE`), with CLI flag taking precedence
- `NewNetBird()` accepts the setup key file path at construction time
- systemd unit file documents the optional `NETBIRD_SETUP_KEY_FILE` environment variable

## Capabilities

### New Capabilities

- `netbird-setup-key`: Configuration and passing of `--setup-key-file` to `netbird up` when a path is provided

### Modified Capabilities

- `netbird-watcher`: The reconnect behavior now conditionally includes `--setup-key-file` in the `netbird up` invocation

## Impact

- `watcher/netbird.go`: `NetBird` struct gains a `setupKeyFile` field; `NewNetBird` accepts a path argument
- `main.go`: Adds flag parsing for `--setup-key-file` and reads `NETBIRD_SETUP_KEY_FILE` env var
- `nestbird.service`: Documents optional `NETBIRD_SETUP_KEY_FILE` environment variable
- `NetBirdClient` interface: unchanged
