package kubenet

import (
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

func NewLinkModels(bandwidth, singleHopLatency string) ([]kubenet_config.LinkModel, error) {
	linkModels := []kubenet_config.LinkModel{}
	satelliteLaserLinkParam := kubenet_config.SatelliteLaserLinkParam{
		Type:             "SatelliteLaserLink",
		Bandwidth:        bandwidth,
		SingleHopLatency: singleHopLatency,
	}
	linkModelParam := kubenet_config.LinkModelParam{
		SatelliteLaserLinkParam: &satelliteLaserLinkParam,
	}
	linkModel := kubenet_config.LinkModel{
		Name:  "linkmodel1",
		Param: linkModelParam,
	}
	linkModels = append(linkModels, linkModel)
	return linkModels, nil
}

func NewSatelliteLaserLink(bandwidth, singleHopLatency string) kubenet_config.LinkModel {
	satelliteLaserLinkParam := kubenet_config.SatelliteLaserLinkParam{
		Type:             "SatelliteLaserLink",
		Bandwidth:        bandwidth,
		SingleHopLatency: singleHopLatency,
	}
	linkModelParam := kubenet_config.LinkModelParam{
		SatelliteLaserLinkParam: &satelliteLaserLinkParam,
	}
	linkModel := kubenet_config.LinkModel{
		Name:  "satlaser",
		Param: linkModelParam,
	}
	return linkModel
}

func NewGeneralWirelessLink() kubenet_config.LinkModel {
	generalWirelessLinkParam := kubenet_config.GeneralWirelessLinkParam{
		Type:      "GeneralWirelessLink",
		Bandwidth: 0.0,
	}
	linkModelParam := kubenet_config.LinkModelParam{
		GeneralWirelessLinkParam: &generalWirelessLinkParam,
	}
	linkModel := kubenet_config.LinkModel{
		Name:  "wireless",
		Param: linkModelParam,
	}
	return linkModel
}

func MergeLinkModels(lms1, lms2 *[]kubenet_config.LinkModel) {
	for _, lm := range *lms2 {
		if notInLms(lm.Name, *lms1) {
			*lms1 = append(*lms1, lm)
		}
	}
}

func notInLms(name string, lms []kubenet_config.LinkModel) bool {
	for _, lm := range lms {
		if name == lm.Name {
			return false
		}
	}
	return true
}
