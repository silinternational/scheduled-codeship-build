#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# allow error checking test to pass
unset CS_BUILD_REFERENCE

go test -v ./cron/builder/
