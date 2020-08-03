#!/usr/bin/env bash

set -e

ROOT="$(readlink -f $(dirname ${BASH_SOURCE})/../)"
CODE_ROOT="$ROOT/installer"

BIN_DIR="$ROOT/bin"
BUILD_TARGET="$BIN_DIR/installer"

IMAGE_FILE="$ROOT/build/Dockerfile"
IMAGE_TAG="quay.io/jcope/kni-install"

go_build(){
  (
    set -x
    cd "$CODE_ROOT"
    echo "compiling package: $(pwd)/..."
    export GOOS=linux GOARCH=amd64
    go build -o "$BIN_DIR/" ./...
  )
}

docker_build(){
  (
    # set -x
    echo "copying $IMAGE_FILE => $BUILD_DIR"
    cp "$IMAGE_FILE" "$BUILD_DIR"
    echo "copying $BUILD_TARGET => $BUILD_DIR"
    cp "$BUILD_TARGET" "$BUILD_DIR"
    docker build -t "$IMAGE_TAG" "$BUILD_DIR"
  )
}

main(){
  go_build
  BUILD_DIR="$(mktemp -d)"
  docker_build
  rm -rf "$BUILD_DIR"
}

main