#!/bin/bash
#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core build
#

set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${PROJECT_ROOT}/build/common.sh"

build_docker_image core/build build-server
