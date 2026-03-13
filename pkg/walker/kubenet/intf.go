package kubenet

import (
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

// 生成卫星接口的参数结构体
type intfDesc struct {
	Name          string
	TopologyModel string
	LinkModel     string
}

// X XX Y YY方向四个intf
var defaultIntfs = []intfDesc{
	{X_INTF_NAME, INTRA_ORBIT_ISL_TOPOLOGY_MODEL, SAT_LINK_MODEL},
	{XX_INTF_NAME, INTRA_ORBIT_ISL_TOPOLOGY_MODEL, SAT_LINK_MODEL},
	{Y_INTF_NAME, INTER_ORBIT_ISL_TOPOLOGY_MODEL, SAT_LINK_MODEL},
	{YY_INTF_NAME, INTER_ORBIT_ISL_TOPOLOGY_MODEL, SAT_LINK_MODEL},
}

// Z ZZ方向跨壳层intf
var crossSphereIntfs = []intfDesc{
	{Z_INTF_NAME, INTER_SPHERE_ISL_TOPOLOGY_MODEL, SAT_LINK_MODEL},
	{ZZ_INTF_NAME, INTER_SPHERE_ISL_TOPOLOGY_MODEL, SAT_LINK_MODEL},
}

type IntfGenerateParam struct {
	RoutingAreaId int32
	CrossSphere   bool
	GroundAccess  bool
}

func (g *KubenetGenerator) GenerateSatelliteIntfs(param IntfGenerateParam) ([]kubenet_config.Intf, error) {
	intfs := make([]kubenet_config.Intf, 0, len(defaultIntfs))

	var routingAreaPtr *int32
	if param.RoutingAreaId > 0 {
		routingAreaPtr = &param.RoutingAreaId
	}

	// X XX Y YY方向四个intf
	for _, d := range defaultIntfs {
		intfs = append(intfs, kubenet_config.Intf{
			Name:          d.Name,
			TopologyModel: d.TopologyModel,
			LinkModel:     d.LinkModel,
			TopologyParam: *g.IntfParamSet.WalkerInterSatTopoParams[d.Name],
			LinkParam:     *g.IntfParamSet.InterSatLinkParam,
			RoutingArea:   routingAreaPtr,
		})
	}
	// Z ZZ方向跨壳层intf
	if param.CrossSphere {
		for _, d := range crossSphereIntfs {
			intfs = append(intfs, kubenet_config.Intf{
				Name:          d.Name,
				TopologyModel: d.TopologyModel,
				LinkModel:     d.LinkModel,
				TopologyParam: *g.IntfParamSet.WalkerInterSatTopoParams[d.Name],
				LinkParam:     *g.IntfParamSet.InterSatLinkParam,
				RoutingArea:   routingAreaPtr,
			})
		}
	}

	// 地面接入接口
	if param.GroundAccess {
		accessIntf := kubenet_config.Intf{
			Name:          ACCESS_INTF_NAME,
			TopologyModel: g.Config.AccessTopologyModel,
			LinkModel:     g.Config.AccessLinkModel,
			TopologyParam: *g.IntfParamSet.ConeTopologyParam,
			LinkParam:     *g.IntfParamSet.AccessIntfLinkParam,
			RoutingArea:   routingAreaPtr,
		}
		intfs = append(intfs, accessIntf)
	}

	return intfs, nil
}
