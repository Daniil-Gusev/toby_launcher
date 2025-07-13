#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <version> (e.g., $0 v1.0.0)"
  exit 1
fi

VERSION="$1"
APP_NAME="TobyLauncher"
PROJECT_NAME="toby_launcher"
APP_MODULE_PATH="${PROJECT_NAME}/core/version"
INSTALLER_MODULE_PATH="main"
PLATFORMS=("linux/amd64" "linux/arm64" "windows/amd64" "windows/arm64" "darwin/amd64" "darwin/arm64")
RELEASE_DIR="release/installers"

if ! command -v go &> /dev/null; then
  echo "go not found, terminate!"
  exit 1
fi

if ! command -v 7z &> /dev/null; then
  echo "7z not found, terminate!"
  exit 1
fi

if [ ! -d "./resources/data" ] || [ -z "$(ls -A ./resources/data)" ]; then
  echo "The directory resources/data must contain files."
  exit 1
fi

mkdir -p tmp
echo "Creating archive tmp/data.7z from resources/data..."
cd resources/data
7z a -mmt=on -mx=6 ../../tmp/data.7z * > /dev/null
cd ../../

./download_gzdoom.sh
if [ $? -ne 0 ]; then
    exit 1
fi

rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"
for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%%/*}
  GOARCH=${PLATFORM##*/}
  APP_BINARY_NAME="$APP_NAME"
  INSTALLER_BINARY_NAME="${APP_NAME}Installer"
  if [ "$GOOS" = "windows" ]; then
    APP_BINARY_NAME="${APP_BINARY_NAME}.exe"
    INSTALLER_BINARY_NAME="${INSTALLER_BINARY_NAME}.exe"
  fi
  rm -rf installer/install/*
  mkdir -p installer/install
  cd $PROJECT_NAME
  echo "Building app binary for $GOOS/$GOARCH..."
  go mod tidy
  BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  CGO=$(  GOOS=$GOOS GOARCH=$GOARCH go env CGO_ENABLED)
  if [ "$(uname)" = "Darwin" ] && [ "$GOOS" = "darwin" ]; then
	  CGO=1
  fi
  GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=$CGO go build -ldflags "-s -w -X ${APP_MODULE_PATH}.AppName=${APP_NAME} -X ${APP_MODULE_PATH}.Version=${VERSION} -X ${APP_MODULE_PATH}.BuildTime=${BUILD_TIME}" -o "../installer/install/${APP_BINARY_NAME}"
  cd ../
  rm -rf tmp/gzdoom*
  if [ "$GOOS" = "linux" ]; then
    GZDOOM_PATH="./resources/gzdoom/linux_${GOARCH}"
  else
    GZDOOM_PATH="./resources/gzdoom/${GOOS}"
  fi
  echo "Creating archive gzdoom.7z from ${GZDOOM_PATH}..."
  mkdir -p tmp/gzdoom
  cp -r $GZDOOM_PATH/* tmp/gzdoom
  cd tmp
  7z a -mmt=on -mx=5 gzdoom.7z gzdoom/* > /dev/null
  rm -f ../installer/data.7z
  echo "Creating full archive installer/data.7z..."
  7z a -mmt=on -mx=0 ../installer/data.7z data.7z gzdoom.7z > /dev/null
  cd ../

  LIB_PATH="./resources/lib/${GOOS}_${GOARCH}"
  if [ -d "$LIB_PATH" ] && [ -n "$(ls -A "$LIB_PATH")" ]; then
	  echo "Copying libraries from ${LIB_PATH} to installer/install/lib..."
    mkdir -p installer/install/lib
    rm -rf installer/install/lib/*
    cp -r "$LIB_PATH"/* installer/install/lib
  fi

  if [ "$GOOS" = "linux" ] || [ "$GOOS" = "darwin" ]; then
    cd installer
    echo "Generating game wrapper for $GOOS..."
    ./generate_wrapper.sh "$GOOS" "$APP_NAME" "$VERSION" "../installer/install/$APP_BINARY_NAME" "../installer/install"
    cd ../
  fi
  if [ "$GOOS" = "darwin" ]; then
    rm -f installer/install/$APP_BINARY_NAME
  fi

  OUTPUT_DIR="${RELEASE_DIR}/${PROJECT_NAME}_installer_${GOOS}_${GOARCH}_${VERSION}"
  mkdir -p "$OUTPUT_DIR"

  cd installer
  echo "Building installer for $GOOS/$GOARCH..."
  go mod tidy
  INSTALLER_OUTPUT="../${OUTPUT_DIR}/${INSTALLER_BINARY_NAME}"
  GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X ${INSTALLER_MODULE_PATH}.AppName=${APP_NAME} -X ${INSTALLER_MODULE_PATH}.BinaryName=${APP_BINARY_NAME}" -o "$INSTALLER_OUTPUT"
  rm -rf data.7z
  rm -rf install/*
  cd ../

  if [ "$GOOS" = "darwin" ]; then
    cd installer
    echo "Generating installer wrapper for $GOOS..."
    ./generate_wrapper.sh "$GOOS" "${APP_NAME}Installer" "$VERSION" "$INSTALLER_OUTPUT" "../$OUTPUT_DIR"
    rm -f "$INSTALLER_OUTPUT"
    cd ../
  fi

  cd "$RELEASE_DIR"
  ARCHIVE_NAME="${PROJECT_NAME}_installer_${GOOS}_${GOARCH}_${VERSION}"
  echo "Creating release tar archive..."
  tar -cf "${ARCHIVE_NAME}.tar" "$ARCHIVE_NAME"
  rm -rf "$ARCHIVE_NAME"
  cd - > /dev/null
done

rm -rf tmp
echo "Installers built in ${RELEASE_DIR}/ directory."