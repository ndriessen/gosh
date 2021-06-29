#!/bin/bash
set -euo pipefail

REQUIRED_COMMANDS=(jq curl)
ALL_COMMANDS_INSTALLED=1
for cmd in "${REQUIRED_COMMANDS[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "Command \"$cmd\" is required to run the install script" >&2
    ALL_COMMANDS_INSTALLED=0
  fi
done

if [ $ALL_COMMANDS_INSTALLED -eq 0 ]; then
  exit 1
fi

PREFIX=${PREFIX:-/usr/local}
INSTALL_BIN=$PREFIX/bin/gosh
GITHUB_REPO_NAME=ndriessen/gosh
GOSH_LATEST_RELEASE_INFO=$(curl --silent --show-error https://api.github.com/repos/$GITHUB_REPO_NAME/releases/latest)
DOWNLOAD_URL=$(jq -r '.assets[]| select(.name == "gosh-linux-amd64").browser_download_url' <<<"$GOSH_LATEST_RELEASE_INFO")
VERIFY_URL=$(jq -r '.assets[]| select(.name == "gosh-linux-amd64.sha256").browser_download_url' <<<"$GOSH_LATEST_RELEASE_INFO")

echo "Downloading and verifying latest gosh binary"
TMP_DOWNLOAD_DIR=$(mktemp -d -t gosh-download-XXXXXXXXXX)
trap 'rm -rf -- "$TMP_DOWNLOAD_DIR"' EXIT

(
  cd "$TMP_DOWNLOAD_DIR" &&
    curl --location --show-error --silent --remote-name "$DOWNLOAD_URL" &&
    curl --location --show-error --silent --remote-name "$VERIFY_URL" &&
    shasum -a 256 -c gosh-linux-amd64.sha256
)

mkdir -p "$(dirname "$INSTALL_BIN")"
install -m0755 "$TMP_DOWNLOAD_DIR/gosh-linux-amd64" "$INSTALL_BIN"

