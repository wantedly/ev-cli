# ev-cli

[![Build Status](https://travis-ci.com/wantedly/ev-cli.svg?token=zsMaScwD2c3BqH37pucy&branch=master)](https://travis-ci.com/wantedly/ev-cli)

CLI tool for managing evaluation data.

## Installation

### For OSX (Using Homebrew)

Formula is avaliable at wantedly/homebrew-tools.

```sh-session
$ brew tap wantedly/tools git@github.com:wantedly/homebrew-tools
$ brew install ev-cli
```

### Other platforms (Linux)

You can download executable file at [HERE](https://github.com/wantedly/ev-cli/releases).

Or, you can download a script for installation.

```console
bash <(curl -sL https://get.wantedlyapp.com/ev-cli)
```

## Prerequisites

ev-cli requires:

- AWS credentials

### AWS credentials

AWS credentials must be set in environment variables or credentials file.

#### Environment variables

```sh-session
$ export AWS_ACCESS_KEY_ID=
$ export AWS_SECRET_ACCESS_KEY=
```

(We strongly recommend to use [envchain](https://github.com/sorah/envchain)).

#### Credential file

Create `~/.aws/credentials` file and write AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY like below:

```
[default]
aws_access_key_id = AK................
aws_secret_access_key = zX..........................
```

`ev-cli` automatically set AWS_REGION as `ap-northeast-1`.

## Usage

- [ev-cli の使い方](https://github.com/wantedly/ev-cli/blob/master/docs/usage.md)

```sh-session
$ ev
CLI tool for managing evaluation data

Usage:
  ev [command]

Available Commands:
  download    Download a file in a target or branch
  export      Export evaluation result files in a target to bigquery
  help        Help about any command
  ls          List targets in a namespace
  ls-branch   List branches in a namespace
  ls-files    List files in a target or branch
  namespaces  List namespaces
  upload      Upload evaluation result files as a target and export it to bigquery
  version     Print the version number of ev

Flags:
  -h, --help   help for ev

Use "ev [command] --help" for more information about a command.
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/wantedly/ev-cli.

## Development

```sh-session
$ make
```
