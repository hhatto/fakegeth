package main

import (
	"github.com/BurntSushi/toml"
)

type HTTPInfo struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type WebsocketInfo struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type IPCInfo struct {
	Path string `toml:"path"`
}

type Config struct {
	Http      *HTTPInfo
	Websocket *WebsocketInfo
	Ipc       *IPCInfo
}

func LoadConfig(path string) (*Config, error) {
	var conf *Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
