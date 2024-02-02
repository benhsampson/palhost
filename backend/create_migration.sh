#!/bin/bash

set -e

migrate create -ext sql -dir services/migrations -seq $1
