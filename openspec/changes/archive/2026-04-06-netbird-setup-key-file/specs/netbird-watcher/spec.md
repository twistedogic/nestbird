## MODIFIED Requirements

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
