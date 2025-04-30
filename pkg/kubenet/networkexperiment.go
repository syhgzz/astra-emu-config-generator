package kubenet

import "dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"

func MergeNetworkExperiment(ne1, ne2 *kubenet_config.NetworkExperiment) {
	ne1.RoutingAreas = ne2.RoutingAreas
	ne1.FadingModels = ne2.FadingModels
	MergeLinkModels(&ne1.LinkModels, &ne2.LinkModels)
	MergeTopoModels(&ne1.TopologyModels, &ne2.TopologyModels)
	MergeConstellations(&ne1.Constellations, &ne2.Constellations)

	ne1.NetworkNodes = append(ne1.NetworkNodes, ne2.NetworkNodes...)
}
