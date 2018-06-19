# GO-SuiteSync

A golang cli implementation of useful filehandling features between local filesystem and netsuite filecabinet

[![CircleCI](https://circleci.com/gh/jroehl/go-suitesync/tree/master.svg?style=svg)](https://circleci.com/gh/jroehl/go-suitesync/tree/master)

> :pushpin: Works only for mac and linux and netsuite accounts on 2018.1 with the 2018.1 sdfcli

- [GO-SuiteSync](#go-suitesync)
  - [Initial Setup](#initial-setup)
    - [Activate token based authentication & SuiteCloud Development Framework](#activate-token-based-authentication--suitecloud-development-framework)
    - [Automatically generate credentials](#automatically-generate-credentials)
    - [Manually generate credentials](#manually-generate-credentials)
  - [Setup circle-ci](#setup-circle-ci)
    - [Setup config](#setup-config)
    - [Setup circle-ci context](#setup-circle-ci-context)
    - [Final setup](#final-setup)
  - [Branch deploy](#branch-deploy)
    - [Config changes](#config-changes)
    - [Variable names](#variable-names)
  - [TODO](#todo)

## Initial Setup

### Activate token based authentication & SuiteCloud Development Framework

1.  Go to Setup > Company > Setup Tasks > Enable Features.
2.  Click the SuiteCloud subtab.
3.  Scroll down to the SuiteScript section, and check the following boxes:
4.  Client SuiteScript.
5.  Server SuiteScript. Click I Agree on the SuiteCloud Terms of Service page.
6.  Scroll to SuiteCloud Development Framework section and check `SUITECLOUD DEVELOPMENT FRAMEWORK`
7.  Click Save.

### Automatically generate credentials

> :rotating_light: The cli setup process currently only works for linux and mac

Download suitesync release and initialize (`suitesync init`) and issue a token to be used for authentication.

```bash
./suitesync/.dependencies/sdfcli issuetoken -url system.netsuite.com -email foo@bar.com -account 123456 -role 3
```

Add `./suitesync/.dependencies/.clicache` content to environment variable `NSCONF_CLITOKEN` for use in automated processes.

> :pencil2: Add the environment variables for local use in `.env` file in directory of `suitesync` cli or in [circle-ci context](https://circleci.com/docs/2.0/contexts/)

### Manually generate credentials

> :bulb: If you are on windows machine use this method

- Create application
  1.  Go to Setup > Integration > Integration Management > Manage Integrations > New
  2.  Enter a name for your application.
  3.  Enter a description, if desired.
  4.  The application State is Enabled by default. (the other option available for selection is blocked.)
  5.  Enter a note, if desired.
  6.  Check the token-based authentication box on the authentication subtab
  7.  Take note of the `consumer key` and `consumer secret`
- Create access token
  1.  Go to `Manage Access Tokens` link in you dashboard (settings portlet)
  2.  Click create access token
  3.  Choose name of application you create beforehand
  4.  Take note of `token id` and `token secret`

Add following variables to the environment for use in automated processes:

```bash
NSCONF_CONSUMER_KEY=consumer key
NSCONF_CONSUMER_SECRET=consumer secret
NSCONF_TOKEN_ID=token id
NSCONF_TOKEN_SECRET=token secret
```

> :pencil2: Add the environment variables for local use in `.env` file in directory of `suitesync` cli or in [circle-ci context](https://circleci.com/docs/2.0/contexts/)

## Setup circle-ci

### Setup config

Use [`ci/circleci-example.yml`](ci/circleci-example.yml) as `.circleci/config.yml` file in repository. For standard sync of one local directory to filecabinet directory this configuration file does not have to be changed.

### Setup [circle-ci context](https://circleci.com/docs/2.0/contexts/)

You will need to set following variables in the [circle-ci context](https://circleci.com/docs/2.0/contexts/) _(! not in the `config.yml` file)_

```bash
########################
####### REQUIRED #######
########################

# Netsuite account
NSCONF_ACCOUNT=123456

# E-mail of the user that generated the credentials
NSCONF_EMAIL=mail@foobar.com

# Automatically generated credentials
NSCONF_CLITOKEN=cli token
#### OR ####
# Manually generated credentials
NSCONF_CONSUMER_KEY=consumer key
NSCONF_CONSUMER_SECRET=consumer secret
NSCONF_TOKEN_ID=token id
NSCONF_TOKEN_SECRET=token secret

# Specify the version of the go-suitesync release to be used
RELEASE=0.0.2
# Specify the realm of the account to be used
NSCONF_REALM=system.netsuite.com
# Specify the relative src directory of the sync
RELATIVE_SRC=./
# Specify the absolute destination directory of the sync
ABSOLUTE_DST=/SuiteScripts/directory

########################
####### OPTIONAL #######
########################

# Specify name of hashfile
NSCONF_HASHFILE=foobar.json # defaults to "hashes.json"
# Specify role of user
NSCONF_ROLE=3 # defaults to "3"

# ONLY USE IF ABSOLUTELY NECESSARY
# Specify password of user (if golang GenerateToken function is used)
NSCONF_PASSWORD=?foo!
# Specify a different rootpath for various file handling calculations (change only if you are sure what you're doing)
NSCONF_ROOTPATH=/Rootpath # defaults to /SuiteScripts
```

> :bulb: Refer to ([Automatically generate credentials](#automatically-generate-credentials) or [Manually generate credentials](#manually-generate-credentials)) for how to generate auth credentials

### Final setup

- Push `config.yml` to remote repository
- Add project to circle-ci
- Trigger initial build

## Branch deploy

If you need to set up different branches to deploy to different accounts you can use following workflow:

> :bookmark: You can find the script that is used to set the environment variables according to the branch [here](https://gist.github.com/jroehl/30d8c212babd5414ad921a298bebec87)

### Config changes

To start building another branch the branch name has to be added to `config.yml`. At the end of the script the new branch name has to be added to every occurence of:

```yaml
...
filters:
  branches:
    only:
      - master
      - new-branch-name
...
```

### Variable names

To add different variable names for different accounts to the [circle-ci context](https://circleci.com/docs/2.0/contexts/) use following syntax (prefix the standard variable names with the branch name and concatenate them with "\_"). Every environment variable can be altered in that way

> :rotating_light: If no branch environment variable is found the fallback is the standard variable (for foo_NSCONF_ACCOUNT = NSCONF_ACCOUNT)

```bash
# account is used for master branch
master_NSCONF_ACCOUNT=master-account
# clitoken is used for development branch
development_NSCONF_CLITOKEN=clitoken-development-branch
```

## TODO

- Tests
- Code cleanup
- Create docker image
- Enable use on windows
- Migrate `sdfcli` calls to SuiteTalk implementation
- Migrate delete restlet to SuiteTalk implementation