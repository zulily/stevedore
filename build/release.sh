#!/bin/bash
#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core build
#

PROJECT_ROOT=$(dirname "${BASH_SOURCE}")/..
GCR_PUSH="${PROJECT_ROOT}/build/gcr-push.sh"

set -o errexit
set -o nounset
set -o pipefail

source "${PROJECT_ROOT}/build/common.sh"

# Pushes docker images to GCS via a (GCS-backed) docker registry
#
# Example:
#   images=("core/foobar" "core/hello-world")
#   push_images images[@]
push_images() {
  echo "Pushing version ${GIT_VERSION} of core/builder to docker repository..."
  $GCR_PUSH -t core/builder -v $GIT_VERSION
}

push_images
