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

	// 只在 YAML 有意义的值时覆盖
	if shellCfg.ShellIndex >= 0 {
		p.ConstellationIdx = shellCfg.ShellIndex
	}
	if shellCfg.Height > 0 {
		p.Height = shellCfg.Height
	}
	if shellCfg.SatelliteNum > 0 {
		p.Num = shellCfg.SatelliteNum
	}
	if shellCfg.PlaneNum > 0 {
		p.Plane = shellCfg.PlaneNum
	}
	if shellCfg.NeededPlaneNum > 0 {
		p.NeededPlane = shellCfg.NeededPlaneNum
	}
	if shellCfg.OrbitInclination != 0 {
		p.Inclination = shellCfg.OrbitInclination
	}
	if shellCfg.NameRule != "" {
		p.NameRule = shellCfg.NameRule
	}
	if shellCfg.InitialDelta >= 0 {
		p.DeltaInit = shellCfg.InitialDelta
	}
	p.Phase = shellCfg.Phase

	return nil
}
