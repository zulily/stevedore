#!/bin/bash
#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core stevedore
#
# Runs test for a golang package, then builds a statically-linked binary with no C dependencies or
# debugging info, from a golang package path.  The package must contain a "main".
#
# Usage: build-and-copy.sh git.zulily.com/core_backend/reasoning/apps/ingest
#
# Example: invoking:
#       $ build-and-copy.sh path/to/some/package@
# results in the binary being placed in: $GOPATH/src/path/to/some/package/package

set -o errexit
set -o nounset
set -o pipefail

PKG=$1

# test
echo "testing (with coverage): ${PKG}/..."
godep go test -cover -v $PKG/...

# Get the name of the resulting executable as established by "go build":
#
# The following is basically the bash equivalent of the 'basename' command, but it works on any
# string, not just a path.  Example: "foobar/the/pkg" => "pkg"
echo "outputting to: ./stevedore"

# Build a true, statically-linked binary, with debug info removed, sending the output to the configured path
CGO_ENABLED=0 godep go build -o ./stevedore -a -tags netgo -ldflags -s $PKG
