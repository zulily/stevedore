#!/bin/bash
#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core stevedore
#

set -o errexit
set -o nounset
set -o pipefail

REPO=core-gitlab.corp.zulily.com
NS=core
PROJECT=stevedore
PROJECT_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${PROJECT_ROOT}/build/color.sh"

readonly PREFIX="${REPO}/${NS}/${PROJECT}"

readonly GIT_VERSION=$(git rev-parse HEAD | cut -c 1-8)
readonly DOCKER_ARGS="--no-cache"

# Installs all dependencies for the project (installing godep if not already present),
# runs all tests, then builds/installs all binaries.
init() {
  cd $GOPATH
  go get github.com/tools/godep
  cd $GOPATH/src/$PREFIX
  godep restore
  godep go test ./...
  godep go install ./...
}

# Builds the "base build" docker image (containing the golang runtime, godep, and and all the
# restored Godeps for this repo).
build_base_docker_image() {

  baseImage=$(grep -E "^FROM " ./build/Dockerfile | awk '{print $2}')
  if [ ! -z "$baseImage" ]; then
    echo -e "${GREEN}Pulling base docker container ${baseImage}...${NC}"
    docker pull $baseImage
  fi

  echo -e "${LIGHT_GREEN}Building ${NS}/build...${NC}"
  docker build ${DOCKER_ARGS} -t ${NS}/build ./build
  tag_docker_image ${NS}/build latest
}

# Determines whether or not a docker image with the given name and version exists locally
#
# Example:
#   docker_image_exists "core/reasoning-api", "61b68fa"
docker_image_exists() {
  IMAGE_NAME=$1
  VERSION=$2
  EXISTS=$(docker images | grep ${IMAGE_NAME} | grep ${VERSION} | wc -l)

  if [ "${EXISTS}" == "0" ]; then
    return $(false)
  else
    return $(true)
  fi
}

# Builds a docker image from the supplied target dir, tagging the resulting image name with a
# version.  The supplied target must contain a valid Dockerfile.
#
# Example:
#   build_docker_image "core/reasoning-api", "61b68fa", "/path/to/api"
build_docker_image() {
  local image_name=$1
  local target=$2

  echo -e "${LIGHT_GREEN}Building version ${GIT_VERSION} of ${image_name} from ${target}${NC}"

  # Build the binary, using our "${NS}/build" docker container.  The binary will be placed next to
  # the Dockerfile for the app.
  echo -e "${LIGHT_GREEN}Building golang binary using ${NS}/build container...${NC}"
  docker run --rm -v "$(pwd)":/go/src/${PREFIX} -t ${NS}/build $PREFIX $target

  baseImage=$(grep -E "^FROM " ./Dockerfile | awk '{print $2}')
  if [ ! -z "$baseImage" ]; then
    echo -e "${GREEN}Pulling base docker container ${baseImage}...${NC}"
    docker pull $baseImage
  fi

  echo -e "${LIGHT_GREEN}Building docker container ${image_name}:${GIT_VERSION}...${NC}"
  docker build $DOCKER_ARGS -t $image_name:$GIT_VERSION .
}

# Tags a docker image with the 'latest' tag
#
# Example:
#   tag_docker_image "core/reasoning-api"
tag_docker_image() {
  local image_name=$1
  local version=${2:-$GIT_VERSION}
  echo -e "${LIGHT_GREEN}Tagging ${image_name}, version ${version} as 'latest'${NC}"
  docker tag -f "${image_name}:${version}" "${image_name}:latest"
}

build_go_binary() {
  echo -e "${LIGHT_GREEN}Building and installing target: ${PREFIX}${NC}"
  godep go install -v $PREFIX
  # TODO: test exit codes
}

test_go_project() {
  echo -e "${LIGHT_GREEN}Testing target: ${PREFIX}${NC}"
  godep go test -v $PREFIX
  # TODO: test exit codes
}

