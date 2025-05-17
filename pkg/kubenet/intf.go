package kubenet

import (
	"strconv"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

func NewSatConeIntf(routingArea kubenet_config.RoutingArea) kubenet_config.Intf {
	satelliteConeTopoIntfParam := kubenet_config.SatelliteConeTopoIntfParam{
		Type: "SatelliteConeTopoIntf",
		Role: 1,
	}
	topologyParam := kubenet_config.IntfTopologyParam{
		SatelliteConeTopoIntfParam: &satelliteConeTopoIntfParam,
	}
	generalWirelessLinkIntfParam := kubenet_config.GeneralWirelessLinkIntfParam{
		Type:        "GeneralWirelessLinkIntf",
		FadingModel: "fading_model1",
	}
	linkParam := kubenet_config.IntfLinkParam{
		GeneralWirelessLinkIntfParam: &generalWirelessLinkIntfParam,
	}

	satConeIntf := kubenet_config.Intf{
		Name:          "intf" + strconv.Itoa(4),
		TopologyModel: "cone",
		LinkModel:     "wireless",
		TopologyParam: topologyParam,
		LinkParam:     linkParam,
		RoutingArea:   &routingArea.Id,
	}
	return satConeIntf
}

func NewGroundConeIntf(routingArea kubenet_config.RoutingArea) kubenet_config.Intf {
	satelliteConeTopoIntfParam := kubenet_config.SatelliteConeTopoIntfParam{
		Type: "SatelliteConeTopoIntf",
		Role: 0,
	}
	topologyParam := kubenet_config.IntfTopologyParam{
		SatelliteConeTopoIntfParam: &satelliteConeTopoIntfParam,
	}
	generalWirelessLinkIntfParam := kubenet_config.GeneralWirelessLinkIntfParam{
		Type:        "GeneralWirelessLinkIntf",
		FadingModel: "fading_model1",
	}
	linkParam := kubenet_config.IntfLinkParam{
		GeneralWirelessLinkIntfParam: &generalWirelessLinkIntfParam,
	}

	satConeIntf := kubenet_config.Intf{
		Name:          "intf" + strconv.Itoa(0),
		TopologyModel: "cone",
		LinkModel:     "wireless",
		TopologyParam: topologyParam,
		LinkParam:     linkParam,
		RoutingArea:   &routingArea.Id,
	}
	return satConeIntf
}
