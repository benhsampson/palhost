#!/bin/bash

set -e

if [ -z "${POSTGRESQL_URL}" ]; then
    echo "POSTGRESQL_URL is not set"
    exit 1
fi

if [ $# -lt 1 ] || ([ "$1" != "up" ] && [ "$1" != "down" ]); then
    echo "Usage: $0 <up|down>"
    exit 1
fi

migrate -database "${POSTGRESQL_URL}" -path db/migrations "$1"
