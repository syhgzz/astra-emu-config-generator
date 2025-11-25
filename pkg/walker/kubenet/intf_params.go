package kubenet

import "dslab.sjtu/kubenet/api/pkg/kubenet_config"

type IntfTemplateSet struct {
	WalkerInterSatTopoParams []kubenet_config.IntfTopologyParam
	InterSatLinkParam        *kubenet_config.IntfLinkParam
	ConeTopologyParam        *kubenet_config.IntfTopologyParam
	AccessIntfLinkParam      *kubenet_config.IntfLinkParam
}

type AccessParams struct {
	AccessIntfNum     int
	AccessFadingModel string
}

func NewIntfTemplateSet(enableAccess bool, accessParams AccessParams) *IntfTemplateSet {
	set := &IntfTemplateSet{
		WalkerInterSatTopoParams: []kubenet_config.IntfTopologyParam{
			{WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "X"}},
			{WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "XX"}},
			{WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "Y"}},
			{WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "YY"}},
			{WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "Z"}},
			{WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "ZZ"}},
		},
		InterSatLinkParam: &kubenet_config.IntfLinkParam{
			SatelliteLaserLinkIntfParam: &kubenet_config.SatelliteLaserLinkIntfParam{
				Type: "SatelliteLaserLinkIntf",
			},
		},
	}

	if enableAccess {
		set.ConeTopologyParam = &kubenet_config.IntfTopologyParam{
			SatelliteConeTopoIntfParam: &kubenet_config.SatelliteConeTopoIntfParam{
				Type:     "SatelliteConeTopoIntf",
				Role:     1,
				MaxLinks: int32(accessParams.AccessIntfNum),
			},
		}
		set.AccessIntfLinkParam = &kubenet_config.IntfLinkParam{
			GeneralWirelessLinkIntfParam: &kubenet_config.GeneralWirelessLinkIntfParam{
				Type:        "GeneralWirelessLinkIntf",
				Bandwidth:   1e10,
				Gain:        20,
				PowerInDbm:  55,
				FadingModel: accessParams.AccessFadingModel,
			},
		}
	}

	return set
}
