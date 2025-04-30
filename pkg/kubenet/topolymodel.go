package kubenet

import "dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"

func NewTopologyModels() ([]kubenet_config.TopologyModel, error) {
	topologyModels := []kubenet_config.TopologyModel{}
	walkerInterSatTopoParam := kubenet_config.WalkerInterSatTopoParam{
		Type: "WalkerInterSatTopo",
	}
	topologyModelParam := kubenet_config.TopologyModelParam{
		WalkerInterSatTopoParam: &walkerInterSatTopoParam,
	}
	topologyModel := kubenet_config.TopologyModel{
		Name:  "topologymodel1",
		Param: topologyModelParam,
	}
	topologyModels = append(topologyModels, topologyModel)
	return topologyModels, nil
}

func NewWalkerInterSatTopoModel() kubenet_config.TopologyModel {
	walkerInterSatTopoParam := kubenet_config.WalkerInterSatTopoParam{
		Type: "WalkerInterSatTopo",
	}
	topologyModelParam := kubenet_config.TopologyModelParam{
		WalkerInterSatTopoParam: &walkerInterSatTopoParam,
	}
	topologyModel := kubenet_config.TopologyModel{
		Name:  "walkerinter",
		Param: topologyModelParam,
	}
	return topologyModel
}

func NewSatelliteConeTopoModel() kubenet_config.TopologyModel {
	satelliteConeTopoParam := kubenet_config.SatelliteConeTopoParam{
		Type: "SatelliteConeTopo",
	}
	topologyModelParam := kubenet_config.TopologyModelParam{
		SatelliteConeTopoParam: &satelliteConeTopoParam,
	}
	topologyModel := kubenet_config.TopologyModel{
		Name:  "cone",
		Param: topologyModelParam,
	}
	return topologyModel
}

func MergeTopoModels(tms1, tms2 *[]kubenet_config.TopologyModel) {
	for _, tm := range *tms2 {
		if notInTms(tm.Name, *tms1) {
			*tms1 = append(*tms1, tm)
		}
	}
}

func notInTms(name string, tms []kubenet_config.TopologyModel) bool {
	for _, tm := range tms {
		if name == tm.Name {
			return false
		}
	}
	return true
}

func MergeConstellations(cs1, cs2 *[]kubenet_config.Constellation) {
	*cs1 = append(*cs1, *cs2...)
}
