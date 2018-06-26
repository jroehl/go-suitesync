#!/usr/bin/env bash

# Bash script to download release from github
#
# Usage in a script:
# $ export RELEASE=0.0.3
# $ export REPO_OWNER=jroehl
# $ export REPO=go-suitesync
# $ export EXECUTABLE=suitesync
# $ curl -sL https://raw.githubusercontent.com/jroehl/go-suitesync/master/ci/download-release.sh | bash

# Define variables.
if [[ "$OSTYPE" == "linux-gnu" ]]; then
  TAR="${REPO}_${RELEASE}_linux_amd64.tar.gz"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  TAR="${REPO}_${RELEASE}_darwin_amd64.tar.gz"
else
  echo "Only \"MacOS\" and \"Linux\" are supported - not \"$OSTYPE\""
  exit 1;
fi

TAG="v$RELEASE"
EXECUTABLE_DIR=$(pwd)/.$EXECUTABLE

GH_API="https://api.github.com"
GH_REPO="$GH_API/repos/${REPO_OWNER}/${REPO}"
GH_TAGS="$GH_REPO/releases/tags/$TAG"

# create dir
mkdir -p $EXECUTABLE_DIR
cd $EXECUTABLE_DIR

# Download release
echo ""
echo "Downloading $EXECUTABLE release \"$TAG\""
echo ""

# Read asset tags.
response=$(curl -s $GH_TAGS)

# Get ID of the asset based on given name.
eval $(echo "$response" | grep -C3 "name.:.\+$TAR" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
[ "$id" ] || { echo "Error: Failed to get asset id, response: $response" | awk 'length($0)<100' >&2; exit 1; }

wget --content-disposition --no-cookie -q --header "Accept: application/octet-stream" "$GH_REPO/releases/assets/$id" --show-progress

# unpack
tar -xf $TAR
rm -f $TAR

# make executable
chmod +x $EXECUTABLE_DIR/$EXECUTABLE

echo ""
# smoke test executable installation
if ! [ -f $EXECUTABLE_DIR/$EXECUTABLE ]; then
  echo "$EXECUTABLE setup failed"
  exit 1
fi

echo "$EXECUTABLE download successful"
echo ""