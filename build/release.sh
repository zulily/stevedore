#!/bin/bash
#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core build
#

PROJECT_ROOT=$(dirname "${BASH_SOURCE}")/..
GCR_PUSH=$(which gcr-push.sh)

set -o errexit
set -o nounset
set -o pipefail

if [ -z $GCR_PUSH ]; then
  echo "gcr-push.sh not found on path, looking for CONF_ROOT"
  if [ -z $CONF_ROOT ]; then
    echo "CONF_ROOT not defined"
    echo "Please add core_conf/bin to your PATH or specify CONF_ROOT"
    exit 1
  else
    GCR_PUSH="${CONF_ROOT}/bin/gcr-push.sh"
  fi
fi

if [ ! -x $GCR_PUSH ]; then
  echo "${GCR_PUSH} not found or not executable!"
  exit 1
fi

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
