#!/bin/bash
CONFIG_FILE="../config.env"

set -a
# shellcheck source=$CONFIG_FILE
source "$CONFIG_FILE"
set +a

# reinitialize swagger. This can be removed once swagger is stable
# swag init --parseDependency --parseInternal -g api/router.go

make build
make run
