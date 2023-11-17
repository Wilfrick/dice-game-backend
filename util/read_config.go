package util

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct { //Examples
	WebsocketHost string // "" / "localhost"
	WebsocketPort int    // 32156 / 12345
	WebsocketUrl  string // "localhost" / "aw808.user.srcf.net"
}

func ReadConfigFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
