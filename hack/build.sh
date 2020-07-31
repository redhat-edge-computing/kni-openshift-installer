#!/usr/bin/env bash

ROOT="$(readlink -f $(dirname ${BASH_SOURCE})/../)"
CODE_ROOT="$ROOT/cluster"
BIN_DIR="$CODE_ROOT/bin"
BUILD_TARGET="$BIN_DIR/cluster"
IMAGE_DIR="$ROOT/build"
IMAGE_TAG="quay.io/jcope/cluster"

go_build(){
  (
    # set -x
    cd "$CODE_ROOT" || exit 1
    echo "compiling package: $(pwd)/..."
    export GOOS=linux GOARCH=amd64
    go build -a -o "$BIN_DIR" ./... || exit 2
  )
}

docker_build(){
  (
    # set -x
    cd "$BUILD_DIR" || exit 3
    echo "copying $IMAGE_TAG/Dockerfile => $BUILD_DIR"
    cp "$IMAGE_DIR/Dockerfile" "$BUILD_DIR" || exit 4
    echo "copying $BUILD_TARGET => $BUILD_DIR"
    cp "$BUILD_TARGET" "$BUILD_DIR" || exit 5
    docker build -t "$IMAGE_TAG" .
  )
}

main(){
  go_build
  BUILD_DIR="$(mktemp -d)"
  docker_build
  rm -rf "$BUILD_DIR"
}

main