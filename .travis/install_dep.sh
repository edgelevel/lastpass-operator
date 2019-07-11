#!/bin/bash

# see https://golang.github.io/dep/docs/installation.html

echo "[+] Setup dep"

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

which dep
dep version

echo "[-] Setup dep"
