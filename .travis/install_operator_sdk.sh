#!/bin/bash

# see https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md

set -e

export OPERATOR_SDK_VERSION=v1.5.2

echo "[+] Setup operator-sdk"

curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk_linux_amd64
chmod +x operator-sdk_linux_amd64
sudo mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk

which operator-sdk
operator-sdk version

echo "[-] Setup operator-sdk"
