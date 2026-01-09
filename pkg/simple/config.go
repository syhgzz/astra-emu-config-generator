package simple

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type SimpleConfig struct {
	Topology      string   `yaml:"topology"`
	Num           int      `yaml:"num"`
	Bandwidth     string   `yaml:"bandwidth"`
	Latency       string   `yaml:"latency"`
	ScheduleNodes []string `yaml:"schedule-nodes"`
	NodePrefix    string   `yaml:"node-prefix"`
}

func LoadSimpleConfigFromYAML(path string) (*SimpleConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read yaml: %w", err)
	}

	var cfg SimpleConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &cfg, nil
}
