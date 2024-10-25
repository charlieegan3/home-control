package config

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type Group struct {
	Name  string   `yaml:"name"`
	ID    string   `yaml:"id"`
	Plugs []string `yaml:"plugs"`
}

type Plug struct {
	Name     string `yaml:"name"`
	ID       string `yaml:"id"`
	Hostname string `yaml:"hostname"`
	Enabled  bool   `yaml:"enabled"`
}

type Config struct {
	Groups     []Group
	Plugs      []Plug
	WOLTargets []WOLTarget `yaml:"wol_targets"`
}

type WOLTarget struct {
	Name     string `yaml:"name"`
	Hostname string `yaml:"hostname"`
	MAC      string `yaml:"mac"`
	Online   bool
}

func Parse(r io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(r)

	var c Config
	err := decoder.Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %s", err)
	}

	return &c, nil
}
