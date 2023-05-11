#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# Build binaries
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
$DIR/build.sh

export AWS_REGION="${AWS_REGION}"
export CS_PASSWORD="${CS_PASSWORD}"
export CS_USERNAME="${CS_USERNAME}"

$HOME/.serverless/bin/serverless deploy --verbose --stage prod
