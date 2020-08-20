#!/usr/bin/env bash

set -e

ROOT="$(readlink -f $(dirname ${BASH_SOURCE})/../)"
CODE_ROOT="$ROOT/kni-install"
BUILD_TARGET="$ROOT/bin/kni-install"

(
  set -x
  cd "$ROOT"
  go build -o "$BUILD_TARGET" "$CODE_ROOT/main.go"
)