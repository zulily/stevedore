#!/bin/bash

#
# Pushes a tagged docker image to Google Cloud Storage (GCS) by way of the Google Container Registry.
#

PROJECT=${PROJECT-eternal-empire-754}
VERSION="latest"

function show_help {
script_name=`basename $0`
echo "Pushes a tagged docker image to Google Container Registry"
echo
echo "Usage: ${script_name} -t tag [-p project] [-v version]"
echo
echo "  -t tag           : the name of the image to push"
echo "  -p project       : (optional, default is \"${PROJECT}\") the GCE project"
echo "  -v version       : (optional, default is \"latest\") the version of the image to push"
}

if [ -z "$1" ]; then
  show_help
  exit 1
fi

# The stuff below is from http://stackoverflow.com/questions/192249/how-do-i-parse-command-line-arguments-in-bash
OPTIND=1
while getopts "h?b:t:r:v:" opt; do
  case "$opt" in
    h)
      show_help
      exit 0
      ;;
    t) TAG=$OPTARG
      GCR_TAG="gcr.io/${PROJECT//-/_}/${OPTARG//\//-}"
      ;;
    p) PROJECT=$OPTARG
      ;;
    v) VERSION=$OPTARG
      ;;
  esac
done

shift $((OPTIND-1))

[ "$1" = "--" ] && shift

echo "Beginning push of ${TAG} to ${GCR_TAG}..."

echo "Checking for images named ${TAG}..."

# TODO : if they ever provide a way to filter by tag...
TAGS=`docker images $TAG | awk '{print $2 " " $3}' | grep "${VERSION}" | awk '{print $2}'`

if [ -z "${TAGS}" ]; then
  echo "No image found with tag ${TAG}"
  exit 1
fi

TAG_VERSIONED="${TAG}:${VERSION}"
GCR_TAG_VERSIONED="${GCR_TAG}:${VERSION}"
GCR_TAG_LATEST="${GCR_TAG}:latest"
echo "Tagging ${TAG_VERSIONED} as ${GCR_TAG_VERSIONED}"
docker tag -f "${TAG_VERSIONED}" "${GCR_TAG_VERSIONED}"

echo "Pushing ${GCR_TAG_VERSIONED}..."
gcloud preview docker push "${GCR_TAG_VERSIONED}"

if [ "$VERSION" != "latest" ]; then

  echo "Tagging ${TAG_VERSIONED} as ${GCR_TAG_LATEST}"
  docker tag -f "${TAG_VERSIONED}" "${GCR_TAG_LATEST}"

  echo "Pushing ${GCR_TAG_LATEST}..."
  gcloud preview docker push "${GCR_TAG_LATEST}"

fi
