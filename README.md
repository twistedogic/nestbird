# Nestbird

A simple Go service that monitors NetBird connection status and automatically reconnects when disconnected.

## Features

- Monitors NetBird connection via CLI (`netbird status`)
- Automatically reconnects with `netbird up` when disconnected
- Exponential backoff with jitter (1m base, 5m max)
- Graceful shutdown on SIGINT/SIGTERM
- Structured JSON logging to journald

## Requirements

- Go 1.21+
- NetBird CLI (`netbird`) installed and configured

## Building

```bash
# Build for current platform
make build

# Build for Linux (cross-compile)
make build-linux
```

## Installation

### Manual Installation

1. Build the binary:
   ```bash
   make build
   ```

2. Install the binary:
   ```bash
   sudo make install
   ```

3. Install the systemd service:
   ```bash
   sudo make install-service
   ```

4. Enable and start the service:
   ```bash
   sudo systemctl enable --now nestbird
   ```

### Full Installation

```bash
sudo make full-install
sudo systemctl enable --now nestbird
```

## Uninstallation

```bash
sudo make full-uninstall
```

## Configuration

The watcher polls every 5 minutes. To change the interval or backoff parameters, edit `main.go` and rebuild.

| Parameter | Default | Description |
|-----------|---------|-------------|
| Poll interval | 5m | How often to check connection status |
| Base backoff | 1m | Initial retry wait time |
| Max backoff | 5m | Maximum retry wait time |

## Usage

```bash
# Run directly
./nestbird

# View logs
journalctl -u nestbird -f
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover
```
