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

export CS1_PROJECT_UUID="${CS1_PROJECT_UUID}"
export CS1_BUILD_REFERENCE="${CS1_BUILD_REFERENCE}"

export CS2_PROJECT_UUID="${CS2_PROJECT_UUID}"
export CS2_BUILD_REFERENCE="${CS2_BUILD_REFERENCE}"

serverless deploy --verbose --stage prod
