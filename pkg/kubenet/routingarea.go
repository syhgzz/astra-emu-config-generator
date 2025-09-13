package kubenet

import "dslab.sjtu/kubenet/api/pkg/kubenet_config"

func NewRoutingAreas() ([]kubenet_config.RoutingArea, error) {
	routingAreas := []kubenet_config.RoutingArea{}
	OSPFV3Param := kubenet_config.OSPFV3Param{
		Type:         "OSPFV3",
		Area:         "0.0.0.1",
		Supernet:     "FD00:aaaa:aaaa::1",
		SupernetMask: 96,
		SubnetMask:   112,
	}
	routingAreaParam := kubenet_config.RoutingAreaParam{
		OSPFV3Param: &OSPFV3Param,
	}
	routingArea := kubenet_config.RoutingArea{
		Id:       int32(1),
		Protocol: "OSPFV3",
		Param:    routingAreaParam,
	}
	routingAreas = append(routingAreas, routingArea)
	return routingAreas, nil
}
