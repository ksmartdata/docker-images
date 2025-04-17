#!/bin/bash

# Check if the targets in ci.yml jobs.build-images.matrix.target match the directories in images folder

set -e

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Root directory
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
# CI configuration file path
CI_FILE="${ROOT_DIR}/.github/workflows/ci.yml"
# Images directory path
IMAGES_DIR="${ROOT_DIR}/images"

# Check if files and directories exist
if [ ! -f "${CI_FILE}" ]; then
    echo "Error: CI configuration file ${CI_FILE} does not exist"
    exit 1
fi

if [ ! -d "${IMAGES_DIR}" ]; then
    echo "Error: images directory ${IMAGES_DIR} does not exist"
    exit 1
fi

# Extract target list from CI file
# Use grep and sed to extract values from matrix.target
CI_TARGETS=$(grep -A 20 "matrix:" "${CI_FILE}" | grep -A 10 "target:" | grep -o '"[^"]*"' | sed 's/"//g' | sort)

# Get all subdirectory names in the images directory
IMAGE_DIRS=$(find "${IMAGES_DIR}" -maxdepth 1 -mindepth 1 -type d -exec basename {} \; | sort)

# Compare the two lists
if [ "$(echo "$CI_TARGETS" | md5sum)" = "$(echo "$IMAGE_DIRS" | md5sum)" ]; then
    echo "✓ Targets in CI configuration match directories in images folder"
    exit 0
else
    echo "✗ Targets in CI configuration do not match directories in images folder"
    exit 1
fi