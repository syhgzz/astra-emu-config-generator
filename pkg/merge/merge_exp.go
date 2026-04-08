package merge

import "dslab.sjtu/kubenet/api/pkg/kubenet_config"

func MergeNetworkExperiment(nes []*kubenet_config.NetworkExperiment) *kubenet_config.NetworkExperiment {
	linkModelsMap := make(map[string]kubenet_config.LinkModel)
	routingAreasMap := make(map[int32]kubenet_config.RoutingArea)
	networkNodesMap := make(map[string]kubenet_config.NetworkNode)
	spheres := make([]kubenet_config.ConstellationSpheresInner, 0)
	topologyModelsMap := make(map[string]kubenet_config.TopologyModel)
	fadingModelsMap := make(map[string]kubenet_config.FadingModel)
	networkDeviceModelMap := make(map[string]kubenet_config.NetworkDeviceModel)
	aeroModelsMap := make(map[string]kubenet_config.AeroModel)

	newConstellation := kubenet_config.Constellation{
		Name: MERGED_CONSTELLATION_NAME,
	}

	for _, ne := range nes {
		for _, linkModel := range ne.LinkModels {
			linkModelsMap[linkModel.Name] = linkModel
		}
		for _, routingArea := range ne.RoutingAreas {
			routingAreasMap[routingArea.Id] = routingArea
		}
		for _, networkNode := range ne.NetworkNodes {
			networkNodesMap[networkNode.Name] = networkNode
		}
		// 收集所有constellations里的所有sphere
		for _, constellation := range ne.Constellations {
			for _, s := range constellation.Spheres {
				spheres = append(spheres, s)
			}
		}
		for _, topologyModel := range ne.TopologyModels {
			topologyModelsMap[topologyModel.Name] = topologyModel
		}
		for _, fadingModel := range ne.FadingModels {
			fadingModelsMap[fadingModel.Name] = fadingModel
		}
		for _, networkDeviceModel := range ne.NetworkDeviceModels {
			networkDeviceModelMap[networkDeviceModel.Name] = networkDeviceModel
		}
		if ne.PhysicsSimulation != nil {
			for _, aeroModel := range ne.PhysicsSimulation.AeroModels {
				aeroModelsMap[aeroModel.Name] = aeroModel
			}
		}
	}

	// 给收集到的所有shpere设置idx， 以防重复
	for idx, s := range spheres {
		s.WalkerSphere.Idx = int32(idx)
	}

	// 将收集到的spheres放入新的constellation中
	newConstellation.Spheres = spheres

	ne := &kubenet_config.NetworkExperiment{
		PhysicsSimulation: &kubenet_config.PhysicsSimulation{},
	}
	for _, linkModel := range linkModelsMap {
		ne.LinkModels = append(ne.LinkModels, linkModel)
	}
	for _, routingArea := range routingAreasMap {
		ne.RoutingAreas = append(ne.RoutingAreas, routingArea)
	}
	for _, networkNode := range networkNodesMap {
		ne.NetworkNodes = append(ne.NetworkNodes, networkNode)
	}
	for _, topologyModel := range topologyModelsMap {
		ne.TopologyModels = append(ne.TopologyModels, topologyModel)
	}
	for _, fadingModel := range fadingModelsMap {
		ne.FadingModels = append(ne.FadingModels, fadingModel)
	}
	for _, networkDeviceModel := range networkDeviceModelMap {
		ne.NetworkDeviceModels = append(ne.NetworkDeviceModels, networkDeviceModel)
	}
	for _, aeroModel := range aeroModelsMap {
		ne.PhysicsSimulation.AeroModels = append(ne.PhysicsSimulation.AeroModels, aeroModel)
	}
	ne.Constellations = append(ne.Constellations, newConstellation)
	ne.Extra = nes[0].Extra

	return ne
}
