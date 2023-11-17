package util

import (
	"errors"
	"io/fs"
	"os"
	"testing"
)

func Test_callReadConfigFileNullPathError(t *testing.T) {
	var configStruct Config
	var path string
	configStruct, err := ReadConfigFile(path)
	if err == nil {
		t.Fail()
	}
	_ = configStruct
}

func Test_ReadConfigFileExample(t *testing.T) {
	path := "test_example.json"
	_, err := os.Stat(path)
	if !errors.Is(err, fs.ErrNotExist) {
		t.FailNow()
	}
	file, err := os.Create(path)
	if err != nil {
		t.Errorf("Couldn't open the file: %s \n", err.Error())
	}
	tidy := func() {
		file.Close()
		err = os.Remove(path)
		if err != nil {
			t.Errorf("Couldn't remove the file: %s \n", err.Error())
		}
	}
	defer tidy()

	exampleConfig := "{\"WebsocketHost\": \"\", \"WebsocketPort\": 12345, \"WebsocketUrl\": \"localhost\"}"
	file.WriteString(exampleConfig)
	config, err := ReadConfigFile(path)
	if err != nil {
		t.Fail()
	}
	trueConfig := Config{WebsocketHost: "", WebsocketPort: 12345, WebsocketUrl: "localhost"}
	Assert(t, trueConfig == config)

}
