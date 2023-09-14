#!/bin/bash
CONFIG_FILE="../config.env"

set -a
# shellcheck source=$CONFIG_FILE
source "$CONFIG_FILE"
set +a

make build
make run
