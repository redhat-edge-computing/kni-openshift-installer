#!/usr/bin/env bash

set -e

ROOT="$(readlink -f $(dirname ${BASH_SOURCE})/../)"
CODE_ROOT="$ROOT/kni-install"

BIN_DIR="$ROOT/bin"
BUILD_TARGET="$BIN_DIR/kni-install"

IMAGE_FILE="$ROOT/build/Dockerfile"
IMAGE_TAG="localhost/kni-install"

go_build(){
  (
#    set -x
    cd "$CODE_ROOT"
    echo "compiling package: $(pwd)/... to $BUILD_TARGET"
    export GOOS=$GOOS GOARCH=$GOARCH
    go build -o "$BIN_DIR/" ./...
  )
}

docker_build(){
  (
#     set -x
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

  if [[ $BUILD_IMAGE == 1 ]]; then
    docker_build
  fi

  rm -rf "$BUILD_DIR"
}

while [ ${#@} != 0 ]; do
  case $1 in
    "--image"|"-i")
      BUILD_IMAGE=1
      GOOS=linux
      GOARCH=amd64
      shift 1
      ;;
    "")
      shift 1
      break
      ;;
    *)
      echo "unknown arg: $1"
      exit
      ;;
  esac
done

main