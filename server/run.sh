#!/bin/bash
CONFIG_FILE="../config.env"

set -a
# shellcheck source=$DB_FILE
source "$CONFIG_FILE"
set +a

make build
make run
