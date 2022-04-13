#!/bin/bash

# see https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md

export OPERATOR_SDK_VERSION=v0.17.2

echo "[+] Setup operator-sdk"

curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu
chmod +x operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu
sudo mv operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk

which operator-sdk
operator-sdk version

echo "[-] Setup operator-sdk"
