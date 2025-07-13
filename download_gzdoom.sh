#!/bin/bash

GZDOOM_VERSION="latest"
RESOURCES_DIR="resources/gzdoom"
PLATFORMS=("windows" "darwin" "linux_amd64" "linux_arm64")

all_platforms_ready=true
if [ -d "$RESOURCES_DIR" ]; then
  for platform in "${PLATFORMS[@]}"; do
    if [ ! -d "$RESOURCES_DIR/$platform" ] || [ -z "$(ls -A "$RESOURCES_DIR/$platform")" ]; then
      all_platforms_ready=false
      break
    fi
  done
else
  all_platforms_ready=false
fi

if [ "$all_platforms_ready" = true ]; then
  exit 0
fi

echo "Downloading GZDoom packages..."

mkdir -p tmp
API_URL="https://api.github.com/repos/ZDoom/gzdoom/releases/${GZDOOM_VERSION}"

download_and_process() {
  local platform=$1
  local file_pattern=$2
  local output_dir=$3
  local extract_cmd=$4

  DOWNLOAD_URL=$(curl -s "$API_URL" | jq -r ".assets[] | select(.name | test(\"$file_pattern\")) | .browser_download_url")
  if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Failed to find download URL for $platform."
    exit 1
  fi

  echo "Downloading $platform package from $DOWNLOAD_URL..."
  curl -sL -o "tmp/$platform-package" "$DOWNLOAD_URL" > /dev/null
  if [ $? -ne 0 ]; then
    echo "Error: Failed to download $platform package."
    exit 1
  fi

  mkdir -p "$output_dir"
  eval "$extract_cmd"
  if [ $? -ne 0 ]; then
    echo "Error: Failed to extract or copy $platform package."
    exit 1
  fi
}

download_and_process "darwin" "macos\\\.zip$" "resources/gzdoom/darwin" \
  "unzip -qq -o tmp/darwin-package -d tmp/darwin && cp -r tmp/darwin/GZDoom.app resources/gzdoom/darwin"

download_and_process "windows" "windows\\\.zip$" "resources/gzdoom/windows" \
  "unzip -qq -o tmp/windows-package -d tmp/windows && cp -r tmp/windows/* resources/gzdoom/windows/"

download_and_process "linux_amd64" "amd64\\\.deb$" "resources/gzdoom/linux_amd64" \
  "cd tmp && ar x linux_amd64-package && tar -xJf data.tar.xz && cp -r opt/gzdoom/* ../resources/gzdoom/linux_amd64/ && cd .."

download_and_process "linux_arm64" "arm64\\\.deb$" "resources/gzdoom/linux_arm64" \
  "cd tmp && ar x linux_arm64-package && tar -xJf data.tar.xz && cp -r opt/gzdoom/* ../resources/gzdoom/linux_arm64/ && cd .."

rm -rf tmp

echo "GZDoom binaries successfully downloaded."