package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Config struct {
	Backend struct {
		Pca struct {
			Address  uint16
			Filename string
		}
	}
}

func NewConfig(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}
	var c Config
	err = json.Unmarshal(file, &c)
		return &c, errors.Wrap(err, "reading file")
}