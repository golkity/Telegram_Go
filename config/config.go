package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Token string `json:"token"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteVale, _ := io.ReadAll(file)
	var config Config
	err = json.Unmarshal(byteVale, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
