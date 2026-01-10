package kubenet

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/shell_config"
)

type GeneratorConfig struct {
	SatelliteNum        int
	PlaneNum            int
	NeededPlane         int
	ShellIdx            int
	Phase               float64
	Offset              int32
	Fake                bool
	NodeNamePrefix      string
	AccessIntfNum       int
	AccessFadingModel   string
	AccessLinkModel     string
	AccessTopologyModel string
	RoutingArea         int32
	CrossSphere         bool
	GroundAccess        bool
	GenFlow             string
	FlowBandwidth       string
	Image               string
	NetworkDeviceModel  string
	Schedule            string
}

func NewDefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		SatelliteNum:        300,
		PlaneNum:            10,
		NeededPlane:         10,
		ShellIdx:            1,
		NodeNamePrefix:      "",
		AccessIntfNum:       16,
		AccessFadingModel:   "fading_model1",
		AccessLinkModel:     "Satellite",
		AccessTopologyModel: "cone1",
		NetworkDeviceModel:  "ndm1",
		// default bool var here not work
		Fake:          false,
		GenFlow:       "",
		FlowBandwidth: "",
		CrossSphere:   false,
		GroundAccess:  false,
		Schedule:      "",
	}
}

func LoadGeneratorConfigFromYAML(path string, shellIdx int, cfg *GeneratorConfig) error {
	shellCfg, err := shell_config.LoadShellConfigFromYAML(path, shellIdx)
	if err != nil {
		return fmt.Errorf("load and unmarshal yaml file %s fail: %w", path, err)
	}

	// 只在 YAML 有意义的值时覆盖
	if shellCfg.ShellIndex >= 0 {
		cfg.ShellIdx = shellCfg.ShellIndex
	}
	if shellCfg.SatelliteNum > 0 {
		cfg.SatelliteNum = shellCfg.SatelliteNum
	}
	if shellCfg.PlaneNum > 0 {
		cfg.PlaneNum = shellCfg.PlaneNum
	}
	if shellCfg.NeededPlaneNum > 0 {
		cfg.NeededPlane = shellCfg.NeededPlaneNum
	}
	if shellCfg.RoutingArea > 0 {
		cfg.RoutingArea = shellCfg.RoutingArea
	}
	cfg.Offset = shellCfg.Offset
	cfg.Phase = shellCfg.Phase
	cfg.Fake = shellCfg.Fake
	cfg.GenFlow = shellCfg.GenFlow
	cfg.CrossSphere = shellCfg.CrossSphere
	cfg.GroundAccess = shellCfg.GroundAccess
	cfg.Schedule = shellCfg.Schedule
	cfg.FlowBandwidth = shellCfg.FlowBandwidth

	return nil
}
