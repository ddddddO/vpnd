package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Commands []CommandConfig `json:"commands"`
}

type CommandConfig struct {
	Command         string          `json:"command"`
	Option          string          `json:"option"`
	Check           string          `json:"check"`
	ReCommandConfig ReCommandConfig `json:"reCommandConfig"`
}

type ReCommandConfig struct {
	ReCommand string `json:"reCommand"`
	Option    string `json:"option"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Unmarshal(path string) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read file %s: %v\n", "config.json", err)
	}

	err = json.Unmarshal(raw, &c)
	if err != nil {
		log.Fatalf("failed to json unmarshal: %v\n", err)
	}
}
