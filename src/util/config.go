package util

import (
	"encoding/json"
	"os"
)

type Config struct {
	Providers map[string]any `json:"providers"`
	ViewDelay int            `json:"viewDelay"`
	Server    struct {
		Port int `json:"port"`
	} `json:"server"`
}

func LoadConfig() (Config, error) {
	file, err := os.ReadFile("config.json")
	var config Config
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func CastConfig[T any](config any) (T, error) {
	bytes, err := json.Marshal(config)
	var defaultValue T
	if err != nil {
		return defaultValue, err
	}

	var result T
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return defaultValue, err
	}

	return result, nil
}
