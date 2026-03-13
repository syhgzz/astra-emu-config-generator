package kubenet

import "dslab.sjtu/kubenet/api/pkg/kubenet_config"

type IntfParamSet struct {
	WalkerInterSatTopoParams map[string]*kubenet_config.IntfTopologyParam
	InterSatLinkParam        *kubenet_config.IntfLinkParam
	ConeTopologyParam        *kubenet_config.IntfTopologyParam
	AccessIntfLinkParam      *kubenet_config.IntfLinkParam
	ExternalTopologyParam    *kubenet_config.IntfTopologyParam
}

type AccessParams struct {
	AccessIntfMaxLinkNum int
}

func NewIntfParamSet(enableAccess bool, accessParams AccessParams) *IntfParamSet {
	set := &IntfParamSet{
		WalkerInterSatTopoParams: map[string]*kubenet_config.IntfTopologyParam{
			X_INTF_NAME:  {WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "X"}},
			XX_INTF_NAME: {WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "XX"}},
			Y_INTF_NAME:  {WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "Y"}},
			YY_INTF_NAME: {WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "YY"}},
			Z_INTF_NAME:  {WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "Z"}},
			ZZ_INTF_NAME: {WalkerInterSatTopoIntfParam: &kubenet_config.WalkerInterSatTopoIntfParam{Type: "WalkerInterSatTopoIntf", Axis: "ZZ"}},
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
				MaxLinks: int32(accessParams.AccessIntfMaxLinkNum),
			},
		}
		set.AccessIntfLinkParam = &kubenet_config.IntfLinkParam{
			GeneralWirelessLinkIntfParam: &kubenet_config.GeneralWirelessLinkIntfParam{
				Type:        "GeneralWirelessLinkIntf",
				Bandwidth:   1e10,
				Gain:        20,
				PowerInDbm:  55,
				FadingModel: ACCE_FADING_MODEL,
			},
		}
	}

	set.ExternalTopologyParam = &kubenet_config.IntfTopologyParam{
		ExternalTopoIntfParam: &kubenet_config.ExternalTopoIntfParam{
			Type: "ExternalTopoIntf",
		},
	}

	return set
}
