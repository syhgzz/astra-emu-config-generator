package kubenet

import (
	"fmt"
	"strconv"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

type IntfGenerateParam struct {
	RoutingAreaId int32
	CrossSphere   bool
	GroundAccess  bool
}

func (g *KubenetGenerator) GenerateSatelliteIntfs(param IntfGenerateParam) ([]kubenet_config.Intf, error) {
	intfs := []kubenet_config.Intf{}

	var routingAreaPtr *int32
	if param.RoutingAreaId > 0 {
		routingAreaPtr = &param.RoutingAreaId
	}

	// 前四个默认接口
	for i := 0; i < 4; i++ {
		var topologyModel string
		switch i {
		case 0, 1:
			topologyModel = "intra-orbit-isl"
		case 2, 3:
			topologyModel = "inter-orbit-isl"
		}

		intfs = append(intfs, kubenet_config.Intf{
			Name:          fmt.Sprintf("intf%d", i),
			TopologyModel: topologyModel,
			LinkModel:     "satlaser",
			TopologyParam: g.IntfTemplateSet.WalkerInterSatTopoParams[i],
			LinkParam:     *g.IntfTemplateSet.InterSatLinkParam,
			RoutingArea:   routingAreaPtr,
		})
	}
	// 跨壳层接口
	if param.CrossSphere {
		for i := 4; i < 6; i++ {
			intfs = append(intfs, kubenet_config.Intf{
				Name:          fmt.Sprintf("intf%d", i),
				TopologyModel: "intra-sphere-isl",
				LinkModel:     "satlaser",
				TopologyParam: g.IntfTemplateSet.WalkerInterSatTopoParams[i],
				LinkParam:     *g.IntfTemplateSet.InterSatLinkParam,
				RoutingArea:   routingAreaPtr,
			})
		}
	}

	// 地面接入接口
	if param.GroundAccess {
		accessIntf := kubenet_config.Intf{
			Name:          "acce" + strconv.Itoa(4),
			TopologyModel: g.Config.AccessTopologyModel,
			LinkModel:     g.Config.AccessLinkModel,
			TopologyParam: *g.IntfTemplateSet.ConeTopologyParam,
			LinkParam:     *g.IntfTemplateSet.AccessIntfLinkParam,
			RoutingArea:   routingAreaPtr,
		}
		intfs = append(intfs, accessIntf)
	}

	return intfs, nil
}
