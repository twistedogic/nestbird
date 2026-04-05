## 1. Project Setup

- [x] 1.1 Initialize Go module (`go mod init`)
- [x] 1.2 Create project structure (main.go, watcher/, Makefile)
- [x] 1.3 Add dependencies (none expected for core functionality)

## 2. Core Implementation

- [x] 2.1 Implement backoff algorithm with jitter (`watcher/backoff.go`)
- [x] 2.2 Implement NetBird CLI interaction (`watcher/netbird.go`)
- [x] 2.3 Implement watcher loop (`watcher/watcher.go`)
- [x] 2.4 Implement signal handling and graceful shutdown (`main.go`)

## 3. Testing

- [x] 3.1 Write unit tests for backoff algorithm
- [x] 3.2 Write unit tests for connection detection (mock `netbird status` output)
- [x] 3.3 Test end-to-end with real NetBird CLI

## 4. Deployment

- [x] 4.1 Create systemd unit file (`nestbird.service`)
- [x] 4.2 Add build targets to Makefile (build, install, uninstall)
- [x] 4.3 Document installation process

## 5. Verification

- [x] 5.1 Verify watcher detects disconnection
- [x] 5.2 Verify watcher reconnects automatically
- [x] 5.3 Verify backoff behavior under repeated failures
- [x] 5.4 Verify graceful shutdown on SIGTERM
