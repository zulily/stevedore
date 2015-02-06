#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core stevedore
#
# Old-skool build tools.
#
# Targets (see each target for more information):
#   all/build: Build golang code
#   test: 		 Run tests
#   dockerize: Build/test code and build docker images
#   release:   Deploy and push docker images to registry

all: build

# Build the "base" docker image that all subsequent golang docker images are based on
init:
	build/init.sh
.PHONY: init

# Build golang binaries for pre-configured packages/apps.  To modify the list of
# included packages/apps, see ALL_TARGETS and ALL_IMAGES in build/common.sh.

# Example:
#   make
#   make build
build: init
	build/build.sh
.PHONY: build

# Test go projects
#
# Example:
#   make test
test: build
	build/test.sh
.PHONY: test

# Build docker images and tag them with the latest git SHA
#
# Example:
#   make dockerize
dockerize: build test
		build/dockerize.sh
.PHONY: dockerize

# Release docker images to the GCS-backed docker repository
#
# Example:
#   make release
release: dockerize
	build/release.sh
.PHONY: release
