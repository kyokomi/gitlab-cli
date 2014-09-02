gitlab-cli 
==========

[![Build Status](https://drone.io/github.com/kyokomi/gitlab-cli/status.png)](https://drone.io/github.com/kyokomi/gitlab-cli/latest)
[![Coverage Status](https://img.shields.io/coveralls/kyokomi/gitlab-cli.svg)](https://coveralls.io/r/kyokomi/gitlab-cli?branch=master)
[![GoDoc](https://godoc.org/github.com/kyokomi/gitlab-cli?status.svg)](https://godoc.org/github.com/kyokomi/gitlab-cli)

gitlab command line tool golang

## Install ##

```
$ go get git@github.com:kyokomi/gitlab-cli.git
```
## Usage ##

### Init Config

```
$ gitlab-cli init --host https://gitlab.com/ --api-path api/v3/ --token aaaaaaaaaaa
```

- `--host`: gitlab host url
- `--api-path`: gitlab api version path
- `--token`: your access token

### Issue List

#### Target Gitlab Issues
![/gitlab.png](https://dl.dropbox.com/u/49084962/gitlab.png)

```
$ gitlab-cli list
```

### Create Issue

```
$ gitlab-cli add -t title -d hoge -l aaa,bbbb,hoge,tag
```

- `-t`: issue title
- `-d`: issue detail
- `-l`: issue labels (カンマ区切りで複数可)

## Demo

### Issues List

![/gitlab-cli_demo_issue-list.png](https://dl.dropbox.com/u/49084962/gitlab-cli_demo_issue-list.png)

## LICENSE

[MIT](https://github.com/kyokomi/gitlab-cli/blob/master/LICENSE)

## Author

[kyokomi](https://github.com/kyokomi)

