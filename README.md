# ev-cli

[![Build Status](https://travis-ci.com/wantedly/ev-cli.svg?token=zsMaScwD2c3BqH37pucy&branch=master)](https://travis-ci.com/wantedly/ev-cli)

CLI tool for managing evaluation data.

## Installation

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

ev-cli can display a list of offline evaluation data, upload it, and download it. That data is stored in AWS S3.

```sh-session
$ ev
CLI tool for managing evaluation data

Usage:
  ev [command]

Available Commands:
  download    Download a file in a target or branch
  export      (Used only for debugging ev-export) Export evaluation result files in a target to bigquery
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

## Configuration

You can set s3-bucket used by ev-cli in the configuration file. The default value is `wantedly-evaluate`.

```yml
# ~/.ev/config.yml
bucket: some-bucket
```

The configuration file in ~/.ev/config.yml is used by default, but you can also use the `--config` option to specify another configuration file.

```console
$ ev ls --config ./config.yml
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/wantedly/ev-cli.

## Development

```sh-session
$ make
```
