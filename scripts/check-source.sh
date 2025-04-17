#!/bin/bash

set -e

# Initialize error flag
error_found=0

# Find all Dockerfile files
find_command="find . -name Dockerfile -type f"
dockerfile_files=$(eval $find_command)

if [ -z "$dockerfile_files" ]; then
  echo "No Dockerfile found in the repository."
  exit 0
fi

# Check each Dockerfile
for file in $dockerfile_files; do
  echo "Checking file: $file"
  
  # Check if curl or wget is used to download code from master/main branch
  if grep -E '(curl|wget).*([^a-zA-Z]main[^a-zA-Z]|[^a-zA-Z]master[^a-zA-Z])' "$file" > /dev/null; then
    echo "ERROR: $file contains curl/wget commands that reference master/main branch:"
    grep -n -E '(curl|wget).*([^a-zA-Z]main[^a-zA-Z]|[^a-zA-Z]master[^a-zA-Z])' "$file"
    error_found=1
  fi
  
  # Check if git is used to clone code from master/main branch
  if grep -E 'git (clone|pull|checkout).*([^a-zA-Z]main[^a-zA-Z]|[^a-zA-Z]master[^a-zA-Z])' "$file" > /dev/null; then
    echo "ERROR: $file contains git commands that reference master/main branch:"
    grep -n -E 'git (clone|pull|checkout).*([^a-zA-Z]main[^a-zA-Z]|[^a-zA-Z]master[^a-zA-Z])' "$file"
    error_found=1
  fi
done

# Exit with non-zero status if errors are found
if [ $error_found -eq 1 ]; then
  echo "Error: Found Dockerfile(s) with references to master/main branches."
  echo "Please use specific commit hashes, tags or releases instead of master/main branches."
  exit 1
else
  echo "All Dockerfile files passed the check."
  exit 0
fi