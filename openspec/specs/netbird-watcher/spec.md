## Requirements

### Requirement: Poll NetBird connection status
The watcher SHALL poll `netbird status` every 5 minutes to determine connection state.

#### Scenario: Check status succeeds while connected
- **WHEN** `netbird status` is executed and NetBird is connected
- **THEN** watcher logs "Connected" and sleeps for 5 minutes

#### Scenario: Check status succeeds while disconnected
- **WHEN** `netbird status` is executed and NetBird is disconnected
- **THEN** watcher attempts to reconnect

### Requirement: Reconnect on disconnect
When NetBird is disconnected, the watcher SHALL execute `netbird up` to reconnect. When a setup key file path is configured, `netbird up` SHALL be invoked with `--setup-key-file=<path>`.

#### Scenario: Reconnect succeeds without setup key file
- **WHEN** NetBird is disconnected and no setup key file is configured and `netbird up` succeeds
- **THEN** watcher logs "Connected" and returns to polling

#### Scenario: Reconnect succeeds with setup key file
- **WHEN** NetBird is disconnected and a setup key file path is configured and `netbird up --setup-key-file=<path>` succeeds
- **THEN** watcher logs "Connected" and returns to polling

#### Scenario: Reconnect fails
- **WHEN** NetBird is disconnected and `netbird up` fails (with or without setup key file)
- **THEN** watcher waits with exponential backoff before retrying

### Requirement: Exponential backoff with jitter
On connection failure, the watcher SHALL wait using exponential backoff with jitter, starting at 1 minute and capped at 5 minutes.

#### Scenario: Backoff starts at base interval
- **WHEN** first reconnection attempt fails
- **THEN** watcher waits approximately 1 minute (30s - 90s with jitter)

#### Scenario: Backoff doubles on repeated failures
- **WHEN** reconnection attempts fail consecutively
- **THEN** wait time doubles: 1m → 2m → 4m → 5m (capped)

#### Scenario: Backoff caps at maximum
- **WHEN** wait time would exceed 5 minutes
- **THEN** wait time is capped at 5 minutes (2.5m - 5m with jitter)

### Requirement: Graceful shutdown
The watcher SHALL handle SIGINT and SIGTERM signals to shut down cleanly.

#### Scenario: SIGTERM received
- **WHEN** SIGTERM is sent to the process
- **THEN** watcher exits with code 0

#### Scenario: SIGINT received
- **WHEN** SIGINT (Ctrl+C) is sent to the process
- **THEN** watcher exits with code 0

### Requirement: Log connection state changes
The watcher SHALL log all connection state transitions for debugging.

#### Scenario: Connection established
- **WHEN** NetBird becomes connected
- **THEN** watcher logs "Connected" at INFO level

#### Scenario: Connection lost
- **WHEN** NetBird becomes disconnected
- **THEN** watcher logs "Disconnected, attempting reconnect" at WARN level

#### Scenario: Reconnection failed
- **WHEN** `netbird up` command fails
- **THEN** watcher logs "Reconnection failed, retrying in X" at WARN level

#### Scenario: Reconnection succeeded after failure
- **WHEN** `netbird up` succeeds after previous failures
- **THEN** watcher logs "Reconnected" at INFO level
