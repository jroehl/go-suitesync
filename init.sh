#!/usr/bin/env bash

echo ""
echo "Setting up wrapped clis"
echo "(sdfcli, sdfcli-createproject)"
echo ""

function cmd_exists() {
  command -v "$cmd" >/dev/null 2>&1;
}

function assert_installed() {
  for cmd in "$@"; do
    if ! cmd_exists "$cmd"; then
      echo "Command \"$cmd\" is needed but does not exist"
      exit 1;
    fi
  done
}

if [[ "$OSTYPE" == "linux-gnu" ]]; then
  PLATFORM=linux-x64.tar.gz
  JAVA_SUBDIR=""
elif [[ "$OSTYPE" == "darwin"* ]]; then
  PLATFORM=macosx-x64.tar.gz
  JAVA_SUBDIR=/Contents/Home
else
  echo "Only \"MacOS\" and \"Linux\" are supported - not \"$OSTYPE\""
  exit 1;
fi

# check and install required packages if necessary
(cmd_exists wget || ( cmd_exists brew && brew install wget || cmd_exists apt-get && apt-get install wget || cmd_exists yum && yum install wget)) >/dev/null 2>&1
(cmd_exists tar || ( cmd_exists brew && brew install gnu-tar || cmd_exists apt-get && apt-get install tar || cmd_exists yum && yum install tar)) >/dev/null 2>&1
(cmd_exists find || ( cmd_exists brew && brew install findutils || cmd_exists apt-get && apt-get install findutils || cmd_exists yum && yum install findutils)) >/dev/null 2>&1
(cmd_exists sed || ( cmd_exists brew && brew uninstall gnu-sed && brew install gnu-sed --with-default-names || cmd_exists apt-get && apt-get install sed || cmd_exists yum && yum install sed)) >/dev/null 2>&1
(cmd_exists unzip || ( cmd_exists brew && brew install unzip || cmd_exists apt-get && apt-get install unzip || cmd_exists yum && yum install unzip)) >/dev/null 2>&1

# asset that the required commands exist
assert_installed wget tar find sed unzip

PARENT_DIR=$(pwd)
DEPS_DIR=$PARENT_DIR/.dependencies

# cleanup dep dir
# rm -rf $DEPS_DIR
mkdir -p $DEPS_DIR
cd $DEPS_DIR

_PROGRESS_OPT=""
# check if --show-progess is supported
if wget --help | grep -q '\--show-progress'; then
  _PROGRESS_OPT="-q --show-progress"
fi

MAJOR_VERSION=8
UPDATE_NUMBER=171
VERSION="${MAJOR_VERSION}u${UPDATE_NUMBER}"
BUILD_NUMBER=11
BASE_URL_JRE="http://download.oracle.com/otn-pub/java/jdk/$VERSION-b$BUILD_NUMBER/512cd62ec5174c3487ac17c61aaa89e8/jre-$VERSION-"

echo "Downloading dependencies"
# download jre
wget -nc --no-check-certificate --no-cookies --header "Cookie: oraclelicense=accept-securebackup-cookie" "$BASE_URL_JRE$PLATFORM" $_PROGRESS_OPT

# download all paths from urls file check content-disposition header for name and skip if file exists
wget -i $PARENT_DIR/urls --content-disposition -nc $_PROGRESS_OPT

for tar in *.tar.gz; do tar -xf $tar; done
for jar in *.jar; do unzip -nq $jar; done

JAVA_DIR=$(find $DEPS_DIR -maxdepth 1 -type d -name 'jre*.*')
JAVA_HOME="$JAVA_DIR$JAVA_SUBDIR"
JAVA_DRIVER="$JAVA_HOME/bin/java"

# smoke test java installation
if ! $JAVA_HOME/bin/java -version >/dev/null 2>&1; then
  echo ""
  echo "JRE setup failed"
  exit 1
fi

MAVEN_DIR=$(find $DEPS_DIR -maxdepth 1 -type d -name '*maven*')
MAVEN_BIN=$MAVEN_DIR/bin/mvn

# rewrite sdfcli script
sed -i -e "s|/webdev/sdf/sdk/|$DEPS_DIR|" $DEPS_DIR/sdfcli
sed -i -e "s|mvn|JAVA_HOME=$JAVA_HOME $MAVEN_BIN|" $DEPS_DIR/sdfcli

# create symlinks
rm -f $PARENT_DIR/sdfcli $PARENT_DIR/sdfcli-createproject
ln -s $DEPS_DIR/sdfcli $PARENT_DIR/sdfcli
ln -s $DEPS_DIR/sdfcli-createproject $PARENT_DIR/sdfcli-createproject

# smoke test sdfcli installation and install maven dependencies
if ! ./sdfcli >/dev/null 2>&1; then
  echo ""
  echo "sdfcli setup failed"
  exit 1
fi

echo ""
echo "Setup completed"
echo ""