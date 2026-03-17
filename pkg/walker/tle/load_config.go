package tle

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/shell_config"
)

func LoadTLEParamsFromYAML(path string, shellIdx int, p *TLEParams) error {
	shellCfg, err := shell_config.LoadShellConfigFromYAML(path, shellIdx)
	if err != nil {
		return fmt.Errorf("load and unmarshal yaml file %s fail: %w", path, err)
	}

	if shellCfg.ShellIndex >= 0 {
		p.ConstellationIdx = shellCfg.ShellIndex
	}
	p.Height = shellCfg.Height
	p.Num = shellCfg.SatelliteNum
	p.Plane = shellCfg.PlaneNum
	p.NeededPlane = shellCfg.NeededPlaneNum
	p.NameRule = shellCfg.NameRule
	p.DeltaInit = shellCfg.InitialDelta
	p.Phase = shellCfg.Phase
	p.Inclination = shellCfg.OrbitInclination

	return nil
}
