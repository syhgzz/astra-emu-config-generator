package kubenet

import (
	"dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"
	"dslab.sjtu/kubenet-config-generator/pkg/tle"
)

func NewNetworkNodes(num, plane int, phyConstellation []tle.Satellite,
) ([]kubenet_config.NetworkNode, []kubenet_config.Constellation, error) {
	networkNodes := []kubenet_config.NetworkNode{}
	constellations := []kubenet_config.Constellation{}

	return networkNodes, constellations, nil
}

func NewGroundConeNetworkNode(routingArea kubenet_config.RoutingArea) kubenet_config.NetworkNode {
	containerParam := kubenet_config.ContainerParam{
		Type: "Container",
	}
	networkNodeParam := kubenet_config.NetworkNodeParam{
		ContainerParam: &containerParam,
	}

	visualOptions := kubenet_config.VisualOption{
		ModelId:  5,
		NodeType: 3,
	}

	intf := NewGroundConeIntf(routingArea)
	intfs := []kubenet_config.Intf{}
	intfs = append(intfs, intf)

	staticMobilityModel := kubenet_config.StaticMobilityModel{
		Type: "Static",
	}
	mobilityModel := kubenet_config.NetworkNodeMobilityModel{
		StaticMobilityModel: &staticMobilityModel,
	}

	networkNode := kubenet_config.NetworkNode{
		Name:          "ground",
		Type:          "Container",
		VisualOptions: &visualOptions,
		Intfs:         intfs,
		MobilityModel: mobilityModel,
		Param:         &networkNodeParam,
	}

	return networkNode
}
