package kubenet

import "dslab.sjtu/kubenet/api/pkg/kubenet_config"

func MergeNetworkExperiment(nes []*kubenet_config.NetworkExperiment) *kubenet_config.NetworkExperiment {
	linkModelsMap := make(map[string]kubenet_config.LinkModel)
	routingAreasMap := make(map[int32]kubenet_config.RoutingArea)
	networkNodesMap := make(map[string]kubenet_config.NetworkNode)
	constellationsMap := make(map[string]kubenet_config.Constellation)
	topologyModelsMap := make(map[string]kubenet_config.TopologyModel)
	fadingModelsMap := make(map[string]kubenet_config.FadingModel)

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
		for _, constellation := range ne.Constellations {
			constellationsMap[constellation.Name] = constellation
		}
		for _, topologyModel := range ne.TopologyModels {
			topologyModelsMap[topologyModel.Name] = topologyModel
		}
		for _, fadingModel := range ne.FadingModels {
			fadingModelsMap[fadingModel.Name] = fadingModel
		}
	}

	ne := &kubenet_config.NetworkExperiment{}
	for _, linkModel := range linkModelsMap {
		ne.LinkModels = append(ne.LinkModels, linkModel)
	}
	for _, routingArea := range routingAreasMap {
		ne.RoutingAreas = append(ne.RoutingAreas, routingArea)
	}
	for _, networkNode := range networkNodesMap {
		ne.NetworkNodes = append(ne.NetworkNodes, networkNode)
	}
	for _, constellation := range constellationsMap {
		ne.Constellations = append(ne.Constellations, constellation)
	}
	for _, topologyModel := range topologyModelsMap {
		ne.TopologyModels = append(ne.TopologyModels, topologyModel)
	}
	for _, fadingModel := range fadingModelsMap {
		ne.FadingModels = append(ne.FadingModels, fadingModel)
	}
	ne.Extra = nes[0].Extra

	return ne
}
