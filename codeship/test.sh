#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# allow error checking test to pass (it expects this env var to be missing)
unset CS_PASSWORD

go test -v ./cron/builder/
