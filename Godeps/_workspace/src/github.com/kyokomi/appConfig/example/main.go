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
