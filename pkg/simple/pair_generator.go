package simple

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet_config_util"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
)

func (g *SimpleGenerator) GeneratePair() error {
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
	for i := 0; i < g.config.Num*2; i++ {
		nodeName := fmt.Sprintf("%s-%d", g.config.NodePrefix, i)
		var pairNodeIdx int
		if i%2 == 0 {
			pairNodeIdx = i + 1
		} else {
			pairNodeIdx = i - 1
		}
		pairNodeName := fmt.Sprintf("%s-%d", g.config.NodePrefix, pairNodeIdx)

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

		intf0 := kubenet_config.Intf{
			Name:          "intf0",
			TopologyModel: topologyModelName,
			TopologyParam: kubenet_config.IntfTopologyParam{
				FixedTopoIntfParam: &kubenet_config.FixedTopoIntfParam{
					Type:     "FixedTopoIntf",
					PeerIntf: fmt.Sprintf("%s/intf0", pairNodeName),
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

		if i%2 == 0 {
			backend := "Iperf"
			// Flow Task
			task := kubenet_config.TaskSequence{
				Uid: "1",
				Sequence: []kubenet_config.Task{
					{
						Param: kubenet_config.TaskParam{
							FlowTaskParam: &kubenet_config.FlowTaskParam{
								Type:      "Flow",
								Dst:       pairNodeName,
								Bandwidth: "",
								Backend:   &backend,
							},
						},
					},
				},
			}
			networkNode.Tasks = append(networkNode.Tasks, task)
		}
		networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, networkNode)

		computeNode := g.config.ScheduleNodes[i%len(g.config.ScheduleNodes)]

		networkExperiment.Extra.SchedulePolicy.ExplicitPolicy.ExplicitPairs = append(
			networkExperiment.Extra.SchedulePolicy.ExplicitPolicy.ExplicitPairs,
			kubenet_config.ExplicitPolicyPair{
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
