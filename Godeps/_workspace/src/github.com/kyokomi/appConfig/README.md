# appConfig

[![Build Status](https://drone.io/github.com/kyokomi/appConfig/status.png)](https://drone.io/github.com/kyokomi/appConfig/latest)
[![Coverage Status](https://img.shields.io/coveralls/kyokomi/appConfig.svg)](https://coveralls.io/r/kyokomi/appConfig?branch=master)

=========

homedir read write to app config file for golang ( Go )

## Example

[source](https://github.com/kyokomi/appConfig/blob/master/example/main.go)

```go
package main

import (
	"github.com/kyokomi/appConfig"
	"bytes"
	"log"
	"fmt"
)

func main() {

	accessToken := "aaaaaaaaaaaaaaaaaaaaa"

	c := appConfig.NewAppConfig("goSampleApp")
	err := c.WriteAppConfig(bytes.NewBufferString(accessToken).Bytes())
	if err != nil {
		log.Fatal(err)
	}

	data, err := c.ReadAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}
```

## Install

```sh
$ go get github.com/kyokomi/appConfig
```

## License

[MIT](https://github.com/kyokomi/appConfig/blob/master/LICENSE)

## Author

[kyokomi](https://github.com/kyokomi)
