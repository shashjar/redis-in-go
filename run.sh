#!/bin/sh

# Used to run program locally.

set -e # Exit early if any commands fail

(
  cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory
  go build -o /tmp/redis-in-go app/*.go
)

exec /tmp/redis-in-go "$@"
