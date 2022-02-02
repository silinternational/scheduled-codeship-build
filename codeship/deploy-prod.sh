#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# Build binaries
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
$DIR/build.sh

# Export env vars
export AWS_REGION="${AWS_REGION}"
export CS_ORGANIZATION="${CS_ORGANIZATION}"
export CS_PASSWORD="${CS_PASSWORD}"
export CS_USERNAME="${CS_USERNAME}"
export CS_PROJECT_UUID="${CS_PROJECT_UUID}"
export CS_BUILD_REFERENCE="${CS_BUILD_REFERENCE}"

serverless deploy --verbose --stage prod
