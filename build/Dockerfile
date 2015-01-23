#
# AUTO GENERATED FROM CORE/BOILERPLATE. IF YOU WISH TO MAKE ANY CHANGES PLEASE PUSH THEM THERE
# AND REBUILD THIS FILE WITH
#
#   $ init.sh core-gitlab.corp.zulily.com core build
#
# core/build - builds golang binaries using a pre-built/configured set of golang dependencies
#
# Usage: docker run --rm -v <src_root>:/go/src/<pkg> -t core/build <pkg>
#
# where :
#   <src_root> is the local path to the cloned repository and
#   <pkg>      is the fully-qualified golang package path to build
#
# Example:
#
#  docker run --rm -v "$(pwd)":/go/src/core-gitlab.corp.zulily.com/core/build -t core/build core-gitlab.corp.zulily.com/core/build
#
FROM       golang:latest
MAINTAINER Core Engineering <core@zulily.com>

# We manage all deps using 'godep'
RUN        go get github.com/tools/godep

# Set up a $GOPATH, restore all necessary deps into the $GOPATH
ENV        GOPATH /go/
ENV        GO_PKG core-gitlab.corp.zulily.com/core/build
WORKDIR    $GOPATH/src/$GO_PKG

# Wrapper script will perform the builds when the container is invoked
ADD          ./test-build-copy.sh /test-build-copy.sh
ENTRYPOINT   ["/test-build-copy.sh"]

CMD ["--help"]
