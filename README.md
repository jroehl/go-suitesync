# GO-SuiteSync

A golang implementation of useful filehandling features between local filesystem and netsuite filecabinet

## Initial Setup

### Create Cli token

```bash
./sdfcli issuetoken -url url.netsuite.com -email foo@bar.com -account 123456 -role 3
```

Add `./dependencies/.clicache` content to environment variable `NSCONF_CLITOKEN` for further use in automated processes