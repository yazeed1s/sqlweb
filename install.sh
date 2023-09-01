#!/bin/bash

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
echo "Downloading binary from ${url}..."
curl -LO ${url}
mkdir _temp
unzip ${os}-${arch}.zip -d _temp/
echo "Moving binary to ${bin_dir}..."
mv _temp/bin/${os}-${arch}/${executable} ${bin_dir}
echo "Cleaning up..."
rm -rf "${os}-${arch}.zip" _temp
# chmod +x "${binary_folder}/${executable}"
echo "Executable downloaded and installed successfully!"
