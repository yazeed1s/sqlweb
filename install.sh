#!/bin/bash

exit_with_error() {
  echo "Error: $1" >&2
  exit 1
}

# get the system's architecture for linux and macOS
get_architecture_unix() {
  if [ "$(uname -m)" == "x86_64" ]; then
	  echo "amd-64-bit"
  elif [ "$(uname -m)" == "arm64" ] || [ "$(uname -m)" == "aarch64" ]; then
    echo "arm-64-bit"
  else
    echo "Unsupported architecture: $(uname -m)" >&2; exit 1
  fi
}

# get the system's architecture based on the OS
if [ "$(uname)" == "Darwin" ]; then
	os="mac"
  	arch=$(get_architecture_unix)
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
	os="linux"
  	arch=$(get_architecture_unix)
else
  	echo "Unsupported OS: $(uname)" >&2; exit 1
fi

github_repo="Yazeed1s/sqlweb"
bin_dir="/usr/local/bin"
version="v0.0.1"
executable="sqlweb"
url="https://github.com/${github_repo}/releases/download/${version}/${os}-${arch}.zip"

# if !error_message=$(curl -LO --show-error --fail ${url} 2>&1); then
#     err=$(echo "${error_message}" | grep -o 'curl:.*$')
#     exit_with_error "${err}"
# fi
echo "Downloading binary from ${url}..."
curl -LO --show-error --fail ${url} || exit_with_error "Failed to download binary"

mkdir _temp
unzip ${os}-${arch}.zip -d _temp/ || exit_with_error "Failed to unzip binary"

echo "Moving binary to ${bin_dir}..."
mv _temp/bin/${os}-${arch}/${executable} ${bin_dir} || exit_with_error "Failed to move binary"

echo "Cleaning up..."
rm -rf "${os}-${arch}.zip" _temp || exit_with_error "Failed to clean up"

# chmod +x "${binary_folder}/${executable}"
echo "Executable downloaded and installed successfully!"