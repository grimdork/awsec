# awsec [![CodeQL](https://github.com/grimdork/awsec/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/grimdork/awsec/actions/workflows/codeql-analysis.yml) [![goreleaser](https://github.com/grimdork/awsec/actions/workflows/release.yml/badge.svg)](https://github.com/grimdork/awsec/actions/workflows/release.yml)
Store secrets in AWS Parameter Store.

## What
This tool treats Amazon Web Services Parameter Store as a repository for secrets. Simple strings, encrypted strings (with KMS keys) and string lists are supported there, and this tool helps set and retrieve them more conveniently.

## Installing

### Install from source
If you have Go v1.17 or later installed:
```
go install github.com/grimdork/awsec@v0.11.1
```

Or just clone the project with `git clone git@github.com:grimdork/awsec.git` if you want to contribute.

### Download packages directly
Download a suitable package from [the official package page](https://github.com/grimdork/awsec/releases).

### Homebrew
If you have Homebrew installed on macOS or Linux:
```
brew tap grimdork/tools
brew install grimdork/tools/awsec
```

## Setup
Everything needed to make `aws-cli` run should already be set up. If your company uses Parameter Store, you probably also use AWS tools. You need two comfiguration files at minimum:
- ~/.aws/config
- ~/.aws/credentials

Check AWS documentation for specifics.

## How

### List secrets
The simplest invocation lists all secrets in your configured AWS account:
```
awsec ls
```

This lists every secret in the configured parameter store.

You can also specify the beginning portion of keys to narrow down the list:
```
awsec ls secrets/internal
```

NOTE:
- Parameter Store requires keys to start with a slash. This tool adds it automatically when missing, where it makes sense.
- Keys your IAM user doesn't have access to may still be listed. You still can't fetch their contents.
- You can create policies to set up path-based permissions, limiting certain paths to be accessible only to some users. For instance, you may have a policy for "/secrets*" and another for "/admin*", and users with access to only one can't create or get values starting with the other path. See AWS documentation on IAM policies and groups for further reference.

### Get a secret
```
awsec get secrets/internal/dbpasswords
````
retrieves a parameter named `/secrets/internal/dbpasswords` from the Parameter Store, provided that you have permission to do so.


### Set a secret
```
awsec set secrets/internal/testpw 123456 -s
````
sets the key `secrets/internal/testpw` to `123456` and flags it as secure, which enables AWS KMS encryption.

You can also set string lists (comma-separated values):
```
awsec set secrets/internal/var-list one,1,two,2 -l
```

This sets four values, which well be presented in pairs when you use `get`. This is useful for small configuration files. Technically it's also usable for password lists, but if you want the maximum security use Securestring and split them up.

The `-d` flag allows you to set a description for a key:
```
awsec set -d "This key is a test." secrets/test "This is the test key's value."
```

Finally, it's also possibly to set a key value from a file:
```
awsec set -f secrets/ssh/prod-web prod-web.pem
````
puts the contents of the file `prod-web.pem` into the key `secrets/ssh/prod-web`.

### Tag a secret
AWS allows keys to have tags in addition to descriptions. Tags are used for many things, including filtering billing information. For example:
```
awsec tag secrets/ssh/prod-web -t customer=internal
```

This command updates the secret `secrets/ssh/prod-web` and sets the tag `customer` to `internal`.

NOTE: Each AWS resource can have a maximum of 50 tags.

### Rename a secret
You can rename a key (sort of) like this:
```
awsec rename secrets/ssh/prod-web secrets/ssh/prod-old-web
````

This copies the contents of `secrets/ssh/prod-web` to a key named `secrets/ssh/prod-old-web` and deletes `secrets/ssh/prod-web`.

NOTE: If removal fails because of lacking permissions, you may end up with a duplicate key. Check policies if this happens.

### Remove a key
```
awsec rm secrets/ssh/prod-web
```
removes the key `secrets/ssh/prod-web`, asking to confirm. Use the `-f` flag to skip the question.
