package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
)

type Config struct {
	HomeKit HomeKitConfig
	Doors   []DoorConfig
}

type DoorConfig struct {
	Name         string
	Manufacturer string
	Model        string
	TopGPIO      int `yaml:"top_gpio"`
	BottomGPIO   int `yaml:"bottom_gpio"`
	RelayGPIO    int `yaml:"relay_gpio"`
}

type HomeKitConfig struct {
	PIN         string
	StoragePath string `yaml:"storage_path"`
}

// loadConfiguration reads YAML data from the specified file name and populates
// a Config object.
func loadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	return c, err
}
