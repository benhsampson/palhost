#!bin/bash

set -ex

if [ -z "$1" ]; then
  echo "Usage: $0 <mod_name>"
  exit 2
fi

MOD_NAME="$1"

mkdir "$MOD_NAME"
cd "$MOD_NAME"
touch "$MOD_NAME.go"
go mod init "palhost/$MOD_NAME"
cat >"$MOD_NAME.go" <<EOF
package main
EOF
go work use .
