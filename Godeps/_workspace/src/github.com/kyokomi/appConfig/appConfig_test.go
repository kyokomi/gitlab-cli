package appConfig

import (
	"bytes"
	"fmt"
	"testing"
	"os"
)

func TestWriteReadAppConfig(t *testing.T) {

	c := NewDefaultAppConfig("GoSandbox")

	data := "aaaaaaaaaaaaaaaaaaaaa"
	if err := c.WriteAppConfig(bytes.NewBufferString(data).Bytes()); err != nil {
		t.Error("WriteAppConfig ", err)
	}

	d, err := c.ReadAppConfig()
	if err != nil {
		t.Error("ReadAppConfig ", err)
	}
	if string(d) != data {
		t.Error("ReadAppConfig missmatch ", err)
	}
	fmt.Println(string(d))

	dirPath, err := createAppConfigDirPath("GoSandbox")
	if err != nil {
		t.Error("createAppConfigDirPath ", err)
	}
	if err := os.RemoveAll(dirPath); err != nil {
		t.Error("RemoveAll ", err)
	}
}
