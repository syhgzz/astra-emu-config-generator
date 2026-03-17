package kubenet

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/shell_config"
)

type GeneratorConfig struct {
	SatelliteNum         int
	PlaneNum             int
	NeededPlane          int
	ShellIdx             int
	Phase                float64
	Offset               int32
	Fake                 bool
	NodeNamePrefix       string
	AccessIntfMaxLinkNum int
	RoutingArea          int32
	CrossSphere          bool
	GroundAccess         bool
	ExternalTopo         bool
	GenFlow              string
	FlowBandwidth        string
	Image                string
	NetworkDeviceModel   string
	Schedule             string
}

func NewDefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		AccessIntfMaxLinkNum: 16,
		NetworkDeviceModel:   "ndm1",
	}
}

func LoadGeneratorConfigFromYAML(path string, shellIdx int, cfg *GeneratorConfig) error {
	shellCfg, err := shell_config.LoadShellConfigFromYAML(path, shellIdx)
	if err != nil {
		return fmt.Errorf("load and unmarshal yaml file %s fail: %w", path, err)
	}

	if shellCfg.ShellIndex >= 0 {
		cfg.ShellIdx = shellCfg.ShellIndex
	}
	cfg.SatelliteNum = shellCfg.SatelliteNum
	cfg.PlaneNum = shellCfg.PlaneNum
	cfg.NeededPlane = shellCfg.NeededPlaneNum
	cfg.RoutingArea = shellCfg.RoutingArea
	cfg.Offset = shellCfg.Offset
	cfg.Phase = shellCfg.Phase
	cfg.Fake = shellCfg.Fake
	cfg.GenFlow = shellCfg.GenFlow
	cfg.CrossSphere = shellCfg.CrossSphere
	cfg.GroundAccess = shellCfg.GroundAccess
	cfg.ExternalTopo = shellCfg.ExternalTopo
	cfg.Schedule = shellCfg.Schedule
	cfg.FlowBandwidth = shellCfg.FlowBandwidth
	cfg.Image = shellCfg.Image

	return nil
}
