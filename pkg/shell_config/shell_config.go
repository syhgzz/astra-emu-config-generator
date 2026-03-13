package shell_config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// 定义与 YAML 对应的结构体
type ShellConfig struct {
	ShellIndex       int     `yaml:"shell-index"`
	SatelliteNum     int     `yaml:"satellite-num"`
	PlaneNum         int     `yaml:"plane-num"`
	NeededPlaneNum   int     `yaml:"needed-plane-num"`
	OrbitInclination float64 `yaml:"orbit-inclination"`
	Height           float64 `yaml:"height"`
	Phase            float64 `yaml:"phase"`
	InitialDelta     float64 `yaml:"initial-delta"`
	Offset           int32   `yaml:"offset"`
	Fake             bool    `yaml:"fake"`
	CrossSphere      bool    `yaml:"cross-sphere"`
	GroundAccess     bool    `yaml:"ground-access"`
	ExternalTopo     bool    `yaml:"external-topo"`
	GenFlow          string  `yaml:"gen-flow"`
	FlowBandwidth    string  `yaml:"flow-bandwidth"`
	RoutingArea      int32   `yaml:"routing-area"`
	NameRule         string  `yaml:"name-rule"`
	Schedule         string  `yaml:"schedule"`
	Image            string  `yaml:"image"`
}

type ShellsConfig struct {
	Shells []ShellConfig `yaml:"shells"`
}

func LoadShellConfigFromYAML(path string, idx int) (*ShellConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read yaml: %w", err)
	}

	var cfg ShellsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	if len(cfg.Shells) == 0 {
		return nil, fmt.Errorf("no shells defined in yaml")
	}

	for _, s := range cfg.Shells {
		if s.ShellIndex == idx {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("shell with index %d not found", idx)
}
