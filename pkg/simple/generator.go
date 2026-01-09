package simple

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet_config_util"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
)

type SimpleGenerator struct {
	config     *SimpleConfig
	outputPath string
}

func NewSimpleGenerator(config *SimpleConfig, outputPath string) *SimpleGenerator {
	return &SimpleGenerator{config: config, outputPath: outputPath}
}

func (g *SimpleGenerator) Generate() error {
	// read template
	var networkExperiment kubenet_config.NetworkExperiment

	// Topology Model
	topologyModelParam := kubenet_config.FixedTopoParam{
		Type: "FixedTopo",
	}
	topologyModelName := "fixed"
	topologyModel := kubenet_config.TopologyModel{
		Name: topologyModelName,
		Param: kubenet_config.TopologyModelParam{
			FixedTopoParam: &topologyModelParam,
		},
	}
	networkExperiment.TopologyModels = append(networkExperiment.TopologyModels, topologyModel)

	// Link Model
	linkModelName := "fixed"
	linkModel := kubenet_config.LinkModel{
		Name: linkModelName,
		Param: kubenet_config.LinkModelParam{
			FixedLinkParam: &kubenet_config.FixedLinkParam{
				Type: "FixedLink",
			},
		},
	}
	if g.config.Bandwidth != "" {
		linkModel.Param.FixedLinkParam.Rate = &g.config.Bandwidth
	}
	if g.config.Latency != "" {
		linkModel.Param.FixedLinkParam.Latency = &g.config.Latency
	}
	networkExperiment.LinkModels = append(networkExperiment.LinkModels, linkModel)

	// Routing Area
	var routingAreaId int32 = 1
	routingArea := kubenet_config.RoutingArea{
		Id:       routingAreaId,
		Protocol: "OSPFV3",
		Param: kubenet_config.RoutingAreaParam{
			OSPFV3Param: &kubenet_config.OSPFV3Param{
				Type:         "OSPFV3",
				Area:         "0.0.0.1",
				SubnetMask:   112,
				Supernet:     "FD00:aaaa:aaaa::1",
				SupernetMask: 96,
			},
		},
	}
	networkExperiment.RoutingAreas = append(networkExperiment.RoutingAreas, routingArea)

	// Extra
	networkExperiment.Extra = &kubenet_config.Extra{
		SchedulePolicy: &kubenet_config.ExtraSchedulePolicy{
			ExplicitPolicy: &kubenet_config.ExplicitPolicy{
				Type: "Explicit",
			},
		},
	}

	// Network Nodes
	for i := 0; i < g.config.Num; i++ {
		nodeName := fmt.Sprintf("%s-%d", g.config.NodePrefix, i)

		mobilityModel := kubenet_config.StaticMobilityModel{
			Type: "Static",
		}

		networkNode := kubenet_config.NetworkNode{
			Name:  nodeName,
			Type:  "Container",
			Intfs: []kubenet_config.Intf{},
			MobilityModel: kubenet_config.NetworkNodeMobilityModel{
				StaticMobilityModel: &mobilityModel,
			},
			OrientationModel: kubenet_config.NetworkNodeOrientationModel{
				EarthCenteredPointingOrientationModel: &kubenet_config.EarthCenteredPointingOrientationModel{
					Type: "EarthCenteredPointing",
				},
			},
		}

		if i != 0 {
			lastNodeName := fmt.Sprintf("%s-%d", g.config.NodePrefix, i-1)
			intf0 := kubenet_config.Intf{
				Name:          "intf0",
				TopologyModel: topologyModelName,
				TopologyParam: kubenet_config.IntfTopologyParam{
					FixedTopoIntfParam: &kubenet_config.FixedTopoIntfParam{
						Type:     "FixedTopoIntf",
						PeerIntf: fmt.Sprintf("%s/intf1", lastNodeName),
					},
				},
				LinkModel: linkModelName,
				LinkParam: kubenet_config.IntfLinkParam{
					FixedLinkIntfParam: &kubenet_config.FixedLinkIntfParam{
						Type: "FixedLinkIntf",
					},
				},
				RoutingArea: &routingAreaId,
			}
			networkNode.Intfs = append(networkNode.Intfs, intf0)
		}

		if i != g.config.Num-1 {
			nextNodeName := fmt.Sprintf("%s-%d", g.config.NodePrefix, i+1)
			intf1 := kubenet_config.Intf{
				Name:          "intf1",
				TopologyModel: topologyModelName,
				TopologyParam: kubenet_config.IntfTopologyParam{
					FixedTopoIntfParam: &kubenet_config.FixedTopoIntfParam{
						Type:     "FixedTopoIntf",
						PeerIntf: fmt.Sprintf("%s/intf0", nextNodeName),
					},
				},
				LinkModel: linkModelName,
				LinkParam: kubenet_config.IntfLinkParam{
					FixedLinkIntfParam: &kubenet_config.FixedLinkIntfParam{
						Type: "FixedLinkIntf",
					},
				},
				RoutingArea: &routingAreaId,
			}
			networkNode.Intfs = append(networkNode.Intfs, intf1)
		}

		networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, networkNode)

		computeNode := g.config.ScheduleNodes[i%len(g.config.ScheduleNodes)]

		networkExperiment.Extra.SchedulePolicy.ExplicitPolicy.ExplicitPairs = append(
			networkExperiment.Extra.SchedulePolicy.ExplicitPolicy.ExplicitPairs,
			kubenet_config.ExplictPolicyPair{
				NetworkNode: nodeName,
				ComputeNode: computeNode,
			},
		)
	}

	// write to output
	err := kubenet_config_util.WriteConfig(g.outputPath, &networkExperiment)
	if err != nil {
		logrus.Fatalf("failed to write networkExperiment: %v", err)
	}
	return nil
}
