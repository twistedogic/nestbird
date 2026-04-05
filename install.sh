#!/usr/bin/env bash
set -euo pipefail

REPO="twistedogic/nestbird"
BINARY="nestbird"
INSTALL_DIR="/usr/local/bin"
SERVICE_NAME="nestbird"
SERVICE_PATH="/etc/systemd/system/${SERVICE_NAME}.service"

# Colours
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()    { echo -e "${GREEN}[nestbird]${NC} $*"; }
warn()    { echo -e "${YELLOW}[nestbird]${NC} $*"; }
error()   { echo -e "${RED}[nestbird]${NC} $*" >&2; exit 1; }

# Require root for install/service steps
need_root() {
  if [ "$(id -u)" -ne 0 ]; then
    error "This script must be run as root (use sudo)."
  fi
}

# Detect OS / arch
detect_target() {
  local os arch

  os="$(uname -s)"
  [ "$os" = "Linux" ] || error "Unsupported OS: $os (only Linux is supported)"

  arch="$(uname -m)"
  case "$arch" in
    x86_64)          TARGET="linux-amd64" ;;
    aarch64|arm64)   TARGET="linux-arm64" ;;
    armv7l)          TARGET="linux-armv7" ;;
    *)               error "Unsupported architecture: $arch" ;;
  esac

  info "Detected target: $TARGET"
}

# Fetch latest release tag from GitHub
latest_version() {
  if command -v curl &>/dev/null; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
      | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\(.*\)".*/\1/'
  elif command -v wget &>/dev/null; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" \
      | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\(.*\)".*/\1/'
  else
    error "curl or wget is required"
  fi
}

# Download binary
download_binary() {
  local version="$1"
  local url="https://github.com/${REPO}/releases/download/${version}/${BINARY}-${TARGET}"
  local tmp
  tmp="$(mktemp)"

  info "Downloading ${BINARY} ${version} (${TARGET})..."

  if command -v curl &>/dev/null; then
    curl -fsSL "$url" -o "$tmp"
  else
    wget -qO "$tmp" "$url"
  fi

  chmod +x "$tmp"
  mv "$tmp" "${INSTALL_DIR}/${BINARY}"
  info "Installed to ${INSTALL_DIR}/${BINARY}"
}

# Install systemd service
install_service() {
  if ! command -v systemctl &>/dev/null; then
    warn "systemd not found — skipping service installation"
    return
  fi

  info "Installing systemd service..."
  cat > "$SERVICE_PATH" <<'SERVICE'
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
# Uncomment and set to pass a setup key file to 'netbird up':
# Environment=NETBIRD_SETUP_KEY_FILE=/etc/netbird/setup.key

[Install]
WantedBy=multi-user.target
SERVICE

  systemctl daemon-reload
  systemctl enable --now "$SERVICE_NAME"
  info "Service enabled and started"
  info "View logs with: journalctl -u ${SERVICE_NAME} -f"
}

main() {
  need_root
  detect_target

  local version
  version="${NESTBIRD_VERSION:-$(latest_version)}"
  [ -n "$version" ] || error "Could not determine latest version"
  info "Version: $version"

  download_binary "$version"
  install_service

  echo
  info "Done. nestbird is running."
}

main "$@"
