package helper

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config is project config read from config.json
type Config struct {
	DB struct {
		Driver string `json:"driver"`
		DNS    string `json:"dns"`
	} `json:"db"`
}

// Global is project global variable
var Global struct {
	Config Config
}

func (c *Config) load(file string) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buffer, c)
	if err != nil {
		panic(err)
	}
}

func init() {
	Global.Config.load(os.Getenv("FUNDTOOLROOT") + "/config.json")
}
