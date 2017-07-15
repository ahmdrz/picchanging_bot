package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Telegram struct {
		Token string `json:"token"`
	} `json:"telegram"`
	Database struct {
		Dialect          string `json:"dialect"`
		ConnectionString string `json:"connection_string"`
	} `json:"database"`
}

func ReadConfiguration() (Config, error) {
	var cfg Config
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
