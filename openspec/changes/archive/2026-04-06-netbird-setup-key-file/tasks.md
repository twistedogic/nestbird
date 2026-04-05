## 1. watcher/netbird.go

- [x] 1.1 Add `setupKeyFile string` field to `NetBird` struct
- [x] 1.2 Update `NewNetBird()` to accept `setupKeyFile string` parameter
- [x] 1.3 Update `Up()` to append `--setup-key-file=<path>` to args when `setupKeyFile` is non-empty

## 2. main.go

- [x] 2.1 Add `--setup-key-file` flag using `flag` package
- [x] 2.2 Fall back to `NETBIRD_SETUP_KEY_FILE` env var when flag is empty
- [x] 2.3 Pass resolved path to `NewNetBird()`

## 3. nestbird.service

- [x] 3.1 Add commented-out `Environment=NETBIRD_SETUP_KEY_FILE=` line with usage note

## 4. Tests

- [x] 4.1 Add test for `Up()` builds correct args when `setupKeyFile` is set
- [x] 4.2 Add test for `Up()` builds correct args when `setupKeyFile` is empty
