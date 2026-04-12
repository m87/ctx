#!/usr/bin/env sh
set -eu

SRC_DIR="${1:-ui/dist/ctx/browser}"
DEST_DIR="${2:-server/spa/dist}"

if [ ! -d "$SRC_DIR" ]; then
  echo "missing SPA build directory: $SRC_DIR" >&2
  echo "run: (cd ui && npm ci && npm run build -- --configuration production)" >&2
  exit 1
fi

mkdir -p "$DEST_DIR"
find "$DEST_DIR" -mindepth 1 -maxdepth 1 ! -name '.keep' -exec rm -rf {} +
cp -a "$SRC_DIR"/. "$DEST_DIR"/
[ -f "$DEST_DIR/.keep" ] || touch "$DEST_DIR/.keep"
