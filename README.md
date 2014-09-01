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
$ gitlab-cli l
[56091(22)] opened : [hgoe] (kyokomi)
[56090(21)] opened : [hgoe] ()
[53679(18)] opened : [たいとるだよー] (kyokomi)
```

### Create Issue

```
$ gitlab-cli i -t title -d hoge -l aaa,bbbb,hoge,tag
```

- `-t`: issue title
- `-d`: issue detail
- `-l`: issue labels (カンマ区切りで複数可)

## LICENSE

MIT

## Author

[kyokomi](https://github.com/kyokomi)

