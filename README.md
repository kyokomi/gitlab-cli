gitlab-cli
==========

gitlab command line tool golang

## Install ##

```
$ go get git@github.com:kyokomi/gitlab-cli.git
```

## Setup ##

### config.json Sample

`$HOME/.gitlab-cli/config.json`

```
{
  "host":     "https://gitlab.com/",
  "api_path": "api/v3/",
  "token":    "aaaaaaaaaaaaaaaaaaaaaaa"
}
```

## Usage ##

```
$ gitlab-cli --help
NAME:
   gitlab-cli - todo:

USAGE:
   gitlab-cli [global options] command [command options] [arguments...]

VERSION:
   0.0.2

COMMANDS:
   create_issue, i  project create issue
   check-project, c check project name
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gitlab.skip-cert-check If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack.
   --version, -v        print the version
   --help, -h           show help

$ gitlab-cli i --help
NAME:
   create_issue - project create issue

USAGE:
   command create_issue [command options] [arguments...]

DESCRIPTION:


OPTIONS:
   --title, -t      issue title.
   --description, -d    issue description.
   --label, -l      label example hoge,fuga,piyo.
```

### Create Issue

```
$ gitlab-cli i -t title -d hoge -l aaa,bbbb,hoge,tag
```

## LICENSE

MIT

