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

// 星地链路接入intf
var acceIntfs = []intfDesc{
	{ACCESS_INTF_NAME, ACCE_TOPOLOGY_MODEL, ACCE_LINK_MODEL},
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

	// 1. 生成默认的X XX Y YY四个方向的intf
	intfs = append(intfs, g.generateIslIntfs(routingAreaPtr, defaultIntfs)...)

	// 2. 生成跨壳层的Z ZZ方向的intf
	if param.CrossSphere {
		intfs = append(intfs, g.generateIslIntfs(routingAreaPtr, crossSphereIntfs)...)
	}

	// 3. 生成地面接入intf
	if param.GroundAccess {
		intfs = append(intfs, g.generateAccessIntfs(routingAreaPtr, acceIntfs)...)
	}

	// 4. 如果是外部拓扑，则修改所有intf的拓扑模型为external
	if g.Config.ExternalTopo {
		g.editExternalTopo(intfs)
	}

	return intfs, nil
}

func (g *KubenetGenerator) buildIntf(
	name string,
	topologyModel string,
	linkModel string,
	topoParam kubenet_config.IntfTopologyParam,
	linkParam kubenet_config.IntfLinkParam,
	routingArea *int32,
) kubenet_config.Intf {
	return kubenet_config.Intf{
		Name:          name,
		TopologyModel: topologyModel,
		LinkModel:     linkModel,
		TopologyParam: topoParam,
		LinkParam:     linkParam,
		RoutingArea:   routingArea,
	}
}

// 生成星间链路intf
func (g *KubenetGenerator) generateIslIntfs(routingArea *int32, desc []intfDesc) []kubenet_config.Intf {
	intfs := []kubenet_config.Intf{}

	for _, d := range desc {
		intfs = append(intfs,
			g.buildIntf(
				d.Name,
				d.TopologyModel,
				d.LinkModel,
				*g.IntfParamSet.WalkerInterSatTopoParams[d.Name],
				*g.IntfParamSet.InterSatLinkParam,
				routingArea,
			),
		)
	}
	return intfs
}

// 生成地面接入intf
func (g *KubenetGenerator) generateAccessIntfs(routingArea *int32, desc []intfDesc) []kubenet_config.Intf {
	intfs := []kubenet_config.Intf{}

	for _, d := range desc {
		intfs = append(intfs,
			g.buildIntf(
				d.Name,
				d.TopologyModel,
				d.LinkModel,
				*g.IntfParamSet.ConeTopologyParam,
				*g.IntfParamSet.AccessIntfLinkParam,
				routingArea,
			),
		)
	}
	return intfs
}

// 生成外部拓扑intf
func (g *KubenetGenerator) editExternalTopo(intfs []kubenet_config.Intf) {
	for i := 0; i < len(intfs); i++ {
		intfs[i].TopologyModel = EXTERNAL_TOPOLOGY_MODEL
		intfs[i].TopologyParam = *g.IntfParamSet.ExternalTopologyParam
	}
}
