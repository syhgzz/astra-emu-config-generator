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
		intfs = append(intfs, kubenet_config.Intf{
			Name:          fmt.Sprintf("intf%d", i),
			TopologyModel: "walkerinter",
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
				TopologyModel: "walkerinter",
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
