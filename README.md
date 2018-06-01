# GO-SuiteSync

A golang implementation of useful filehandling features between local filesystem and netsuite filecabinet

## Initial Setup

### Activate token based authentication & SuiteCloud Development Framework

1. Go to Setup > Company > Setup Tasks > Enable Features.
2. Click the SuiteCloud subtab.
3. Scroll down to the SuiteScript section, and check the following boxes:
4. Client SuiteScript.
5. Server SuiteScript. Click I Agree on the SuiteCloud Terms of Service page.
6. Scroll to SuiteCloud Development Framework	section and check SUITECLOUD DEVELOPMENT FRAMEWORK
7. Click Save.

### Create Cli token

```bash
./sdfcli issuetoken -url url.netsuite.com -email foo@bar.com -account 123456 -role 3
```

Add `./dependencies/.clicache` content to environment variable `NSCONF_CLITOKEN` for further use in automated processes

### Create application

1. Go to Setup > Integration > Integration Management > Manage Integrations > New
2. Enter a name for your application.
3. Enter a description, if desired.
4. The application State is Enabled by default. (the other option available for selection is blocked.)
5. Enter a note, if desired.
6. Check the token-based authentication box on the authentication subtab
7. Take note of the `consumer key` and `consumer secret`

### Create access token

1. Go to `Manage Access Tokens` link in you dashboard (settings portlet)
2. Click create access token
3. Choose name of application you create beforehand
4. Take not of `token id` and `token secret`
