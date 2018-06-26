#!/usr/bin/env bash

# Bash script to set environment variables in circleci script
#
# Usage in a script:
# $ curl -sL https://raw.githubusercontent.com/jroehl/go-suitesync/master/ci/set-env.sh | bash
#
# The env variable master_NSCONF_ACCOUNT will overwrite the NSCONF_ACCOUNT variable - if set!

ACCOUNT=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_ACCOUNT")
EMAIL=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_EMAIL")
PASSWORD=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_PASSWORD")
REALM=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_REALM")
ROLE=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_ROLE")
CLITOKEN=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_CLITOKEN")
CONSUMER_KEY=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_CONSUMER_KEY")
CONSUMER_SECRET=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_CONSUMER_SECRET")
TOKEN_ID=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_TOKEN_ID")
TOKEN_SECRET=$(eval echo "\$${CIRCLE_BRANCH}_NSCONF_TOKEN_SECRET")

echo "export NSCONF_ACCOUNT=${ACCOUNT:-$NSCONF_ACCOUNT}" >> $BASH_ENV
echo "export NSCONF_EMAIL=${EMAIL:-$NSCONF_EMAIL}" >> $BASH_ENV
echo "export NSCONF_REALM=${REALM:-$NSCONF_REALM}" >> $BASH_ENV
echo "export NSCONF_ROLE=${ROLE:-$NSCONF_ROLE}" >> $BASH_ENV
echo "export NSCONF_PASSWORD=${PASSWORD:-$NSCONF_PASSWORD}" >> $BASH_ENV
echo "export NSCONF_CLITOKEN=${CLITOKEN:-$NSCONF_CLITOKEN}" >> $BASH_ENV
echo "export NSCONF_CONSUMER_KEY=${CONSUMER_KEY:-$NSCONF_CONSUMER_KEY}" >> $BASH_ENV
echo "export NSCONF_CONSUMER_SECRET=${CONSUMER_SECRET:-$NSCONF_CONSUMER_SECRET}" >> $BASH_ENV
echo "export NSCONF_TOKEN_ID=${TOKEN_ID:-$NSCONF_TOKEN_ID}" >> $BASH_ENV
echo "export NSCONF_TOKEN_SECRET=${TOKEN_SECRET:-$NSCONF_TOKEN_SECRET}" >> $BASH_ENV