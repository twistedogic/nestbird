## ADDED Requirements

### Requirement: Configure setup key file via CLI flag
The service SHALL accept a `--setup-key-file=<path>` CLI flag to specify the path to the NetBird setup key file.

#### Scenario: Flag provided at startup
- **WHEN** nestbird is started with `--setup-key-file=/etc/netbird/setup.key`
- **THEN** the path `/etc/netbird/setup.key` is used for all subsequent `netbird up` invocations

#### Scenario: Flag not provided
- **WHEN** nestbird is started without `--setup-key-file`
- **THEN** no setup key file path is configured from the flag (env var fallback applies)

### Requirement: Configure setup key file via environment variable
The service SHALL read `NETBIRD_SETUP_KEY_FILE` from the environment as a fallback when the CLI flag is not set.

#### Scenario: Env var set, flag not set
- **WHEN** `NETBIRD_SETUP_KEY_FILE=/etc/netbird/setup.key` is set and `--setup-key-file` flag is not provided
- **THEN** the path from the environment variable is used

#### Scenario: Both flag and env var set
- **WHEN** `--setup-key-file=/flag/path.key` is provided and `NETBIRD_SETUP_KEY_FILE=/env/path.key` is set
- **THEN** the CLI flag value `/flag/path.key` takes precedence

#### Scenario: Neither flag nor env var set
- **WHEN** neither `--setup-key-file` nor `NETBIRD_SETUP_KEY_FILE` is configured
- **THEN** nestbird runs without a setup key file (existing behavior preserved)
