package walker

import (
	"fmt"
	"strconv"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet"
	"dslab.sjtu/kubenet-config-generator/pkg/tle"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// kubenetCmd represents the generate command
var kubenetCmd = &cobra.Command{
	Use:   "kubenet",
	Short: "Generate Kubenet Config",
	Long:  `Generate Kubenet Config.`,
	Run: func(cmd *cobra.Command, args []string) {
		var routingAreaPtr *int32
		var networkDeviceModelptr *string

		fmt.Println("Start generating Kubenet Config")
		num, _ := cmd.Flags().GetInt("num")
		constellationIdx, _ := cmd.Flags().GetInt("constellation-idx")
		plane, _ := cmd.Flags().GetInt("plane")
		neededPlane, _ := cmd.Flags().GetInt("needed-plane")
		nodeNamePre, _ := cmd.Flags().GetString("nodeNamePre")
		template, _ := cmd.Flags().GetString("template")
		inputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")
		accessIntfNum, _ := cmd.Flags().GetInt("accessIntfNum")
		accessFadingmodel, _ := cmd.Flags().GetString("accessFadingmodel")
		accessLinkmodel, _ := cmd.Flags().GetString("accessLinkmodel")
		accessTopologymodel, _ := cmd.Flags().GetString("accessTopologymodel")
		offset, _ := cmd.Flags().GetInt32("offset")
		image, _ := cmd.Flags().GetString("image")
		routingAreaId, err := cmd.Flags().GetInt32("routingArea")
		if err == nil {
			routingAreaPtr = &routingAreaId
		}
		networkDeviceModelName, err := cmd.Flags().GetString("networkDeviceModel")
		if err == nil {
			networkDeviceModelptr = &networkDeviceModelName
		}
		fake, _ := cmd.Flags().GetBool("fake")
		crossSphere, _ := cmd.Flags().GetBool("cross-sphere")
		phase, _ := cmd.Flags().GetFloat64("phase")

		enableGroundAccess := false
		if accessIntfNum != 0 && accessLinkmodel != "" &&
			accessFadingmodel != "" && accessTopologymodel != "" {
			enableGroundAccess = true
		}

		satNumPerTrack := num / plane
		actualNum := satNumPerTrack * neededPlane
		isFull := plane == neededPlane

		networkExperiment, err := kubenet.ReadConfig(template)
		if err != nil {
			logrus.Fatalf("failed to read networkExperiment: %v", err)
		}

		networkNodes := []kubenet_config.NetworkNode{}
		constellations := []kubenet_config.Constellation{}

		phyConstellation, err := tle.NewConstellation(inputPath)
		if err != nil {
			logrus.Fatalf("failed to initialize constellation: %v", err)
		}

		// 星间链路
		walkerInterSatTopoIntfParam0 := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
			Axis: "X",
		}
		interSatTopologyParam0 := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam0,
		}

		walkerInterSatTopoIntfParam1 := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
			Axis: "XX",
		}
		interSatTopologyParam1 := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam1,
		}

		walkerInterSatTopoIntfParam2 := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
			Axis: "Y",
		}
		interSatTopologyParam2 := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam2,
		}

		walkerInterSatTopoIntfParam3 := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
			Axis: "YY",
		}
		interSatTopologyParam3 := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam3,
		}

		walkerInterSatTopoIntfParam4 := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
			Axis: "Z",
		}
		interSatTopologyParam4 := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam4,
		}

		walkerInterSatTopoIntfParam5 := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
			Axis: "ZZ",
		}
		interSatTopologyParam5 := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam5,
		}

		satelliteLaserLinkIntfParam := kubenet_config.SatelliteLaserLinkIntfParam{
			Type: "SatelliteLaserLinkIntf",
		}
		interSatlinkParam := kubenet_config.IntfLinkParam{
			SatelliteLaserLinkIntfParam: &satelliteLaserLinkIntfParam,
		}

		// 地面接入链路
		coneTopoIntfParam := kubenet_config.SatelliteConeTopoIntfParam{
			Type:     "SatelliteConeTopoIntf",
			Role:     1,
			MaxLinks: int32(accessIntfNum),
		}
		coneTopologyParam := kubenet_config.IntfTopologyParam{
			SatelliteConeTopoIntfParam: &coneTopoIntfParam,
		}
		generalWirelessLinkIntfParam := kubenet_config.GeneralWirelessLinkIntfParam{
			Type:        "GeneralWirelessLinkIntf",
			Bandwidth:   1e10,
			Gain:        20,
			PowerInDbm:  55,
			FadingModel: accessFadingmodel,
		}
		accessIntfLinkParam := kubenet_config.IntfLinkParam{
			GeneralWirelessLinkIntfParam: &generalWirelessLinkIntfParam,
		}

		testSuiteImage := "harbor.sail.se.sjtu.edu.cn/utils/test-suite:v0.2"

		ports := []int32{5000, 5001, 5002}
		computationUnits := []kubenet_config.ComputationUnit{}
		computationUnits = append(
			computationUnits,
			kubenet_config.ComputationUnit{
				Name:  "testsuite",
				Image: &testSuiteImage,
				Ports: ports,
			},
		)
		if image != "" {
			computationUnits = append(
				computationUnits,
				kubenet_config.ComputationUnit{
					Image: &image,
				},
			)
		}
		containerParam := kubenet_config.ContainerParam{
			Type:             "Container",
			ComputationUnits: computationUnits,
		}
		networkNodeParam := kubenet_config.NetworkNodeParam{
			ContainerParam: &containerParam,
		}

		visualOptions := kubenet_config.VisualOption{
			ModelId:  "1",
			NodeType: 1,
		}

		// Initialize Constellations and NetworkNodes
		satellites := []kubenet_config.OrbitSatellite{}
		planes := []kubenet_config.OrbitPlane{}
		planeIdx := 0
		innerIdx := 0
		for idx := range actualNum {
			nodeName := nodeNamePre + phyConstellation[idx].Name
			satellites = append(
				satellites,
				kubenet_config.OrbitSatellite{
					Idx:       int32(innerIdx),
					Satellite: nodeName,
				},
			)
			innerIdx += 1
			if (innerIdx%satNumPerTrack == 0) && (innerIdx >= satNumPerTrack) {
				innerIdx -= satNumPerTrack
				planes = append(
					planes,
					kubenet_config.OrbitPlane{
						Idx:        int32(planeIdx),
						Satellites: satellites,
					},
				)
				planeIdx += 1
				satellites = []kubenet_config.OrbitSatellite{}
			}
			satelliteMobilityModel := kubenet_config.SatelliteMobilityModel{
				Type:     "Satellite",
				TleLine1: phyConstellation[idx].Satellite.Line1,
				TleLine2: phyConstellation[idx].Satellite.Line2,
			}
			intfs := []kubenet_config.Intf{}
			intf0 := kubenet_config.Intf{
				Name:          "intf0",
				TopologyModel: "walkerinter",
				LinkModel:     "satlaser",
				TopologyParam: interSatTopologyParam0,
				LinkParam:     interSatlinkParam,
				RoutingArea:   routingAreaPtr,
			}
			intfs = append(
				intfs,
				intf0,
			)

			intf1 := kubenet_config.Intf{
				Name:          "intf1",
				TopologyModel: "walkerinter",
				LinkModel:     "satlaser",
				TopologyParam: interSatTopologyParam1,
				LinkParam:     interSatlinkParam,
				RoutingArea:   routingAreaPtr,
			}
			intfs = append(
				intfs,
				intf1,
			)

			intf2 := kubenet_config.Intf{
				Name:          "intf2",
				TopologyModel: "walkerinter",
				LinkModel:     "satlaser",
				TopologyParam: interSatTopologyParam2,
				LinkParam:     interSatlinkParam,
				RoutingArea:   routingAreaPtr,
			}
			intfs = append(
				intfs,
				intf2,
			)

			intf3 := kubenet_config.Intf{
				Name:          "intf3",
				TopologyModel: "walkerinter",
				LinkModel:     "satlaser",
				TopologyParam: interSatTopologyParam3,
				LinkParam:     interSatlinkParam,
				RoutingArea:   routingAreaPtr,
			}
			intfs = append(
				intfs,
				intf3,
			)

			if crossSphere {
				intf4 := kubenet_config.Intf{
					Name:          "intf4",
					TopologyModel: "walkerinter",
					LinkModel:     "satlaser",
					TopologyParam: interSatTopologyParam4,
					LinkParam:     interSatlinkParam,
					RoutingArea:   routingAreaPtr,
				}
				intfs = append(
					intfs,
					intf4,
				)

				intf5 := kubenet_config.Intf{
					Name:          "intf5",
					TopologyModel: "walkerinter",
					LinkModel:     "satlaser",
					TopologyParam: interSatTopologyParam5,
					LinkParam:     interSatlinkParam,
					RoutingArea:   routingAreaPtr,
				}
				intfs = append(
					intfs,
					intf5,
				)
			}

			if enableGroundAccess {

				intfs = append(
					intfs,
					kubenet_config.Intf{
						Name:          "acce" + strconv.Itoa(4),
						TopologyModel: accessTopologymodel,
						LinkModel:     accessLinkmodel,
						TopologyParam: coneTopologyParam,
						LinkParam:     accessIntfLinkParam,
						RoutingArea:   routingAreaPtr,
					},
				)
			}
			// intfs = append(
			// 	intfs,
			// 	kubenet.NewSatConeIntf(routingArea),
			// )
			mobilityModel := kubenet_config.NetworkNodeMobilityModel{
				SatelliteMobilityModel: &satelliteMobilityModel,
			}

			walkerSatelliteOrientationModel := kubenet_config.WalkerSatelliteOrientationModel{
				Type: "WalkerSatellite",
			}

			orientationModel := kubenet_config.NetworkNodeOrientationModel{
				WalkerSatelliteOrientationModel: &walkerSatelliteOrientationModel,
			}

			networkNodes = append(
				networkNodes,
				kubenet_config.NetworkNode{
					Name:               nodeName,
					Type:               "Container",
					VisualOptions:      &visualOptions,
					Intfs:              intfs,
					MobilityModel:      mobilityModel,
					OrientationModel:   orientationModel,
					Param:              &networkNodeParam,
					NetworkDeviceModel: networkDeviceModelptr,
				},
			)
		}
		walkerSphere := kubenet_config.WalkerSphere{
			Type:   "Walker",
			Planes: planes,
			Phase:  int32(phase),
			Offset: offset,
			IsFull: isFull,
		}
		spheres := []kubenet_config.ConstellationSpheresInner{}
		spheres = append(
			spheres,
			kubenet_config.ConstellationSpheresInner{
				WalkerSphere: &walkerSphere,
			},
		)
		constellations = append(
			constellations,
			kubenet_config.Constellation{
				Name:    fmt.Sprintf("s%d", constellationIdx),
				Spheres: spheres,
			},
		)

		// networkNodes = append(networkNodes, kubenet.NewGroundConeNetworkNode(routingArea))

		// Construct NetworkExperiment and convert it to yaml file
		networkExperiment.Constellations = append(networkExperiment.Constellations, constellations...)
		networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, networkNodes...)

		if fake {
			networkExperiment.Extra.IsFakeTest = &fake
		}

		err = kubenet.WriteConfig(outputPath, networkExperiment)
		if err != nil {
			logrus.Fatalf("failed to write networkExperiment: %v", err)
		}

	},
}

func init() {
	WalkerCmd.AddCommand(kubenetCmd)

	// params for constellation
	kubenetCmd.Flags().Int32("offset", 0, "Offset")
	kubenetCmd.Flags().Bool("cross-sphere", false, "cross-sphere")

	// params for network experiment
	kubenetCmd.Flags().String("nodeNamePre", "", "the prefix of network node name")
	kubenetCmd.Flags().StringP("bandwidth", "b", "1Gibps", "linkmodel Bandwidth")
	kubenetCmd.Flags().String("image", "", "application image")
	kubenetCmd.Flags().String("networkDeviceModel", "ndm1", "network device model name")
	kubenetCmd.Flags().Int32P("routingArea", "r", 1, "routing area Id")
	kubenetCmd.Flags().Bool("fake", false, "Fake")

	// params for ground access intfs
	kubenetCmd.Flags().Int("accessIntfNum", 0, "the number of access intfs")
	kubenetCmd.Flags().String("accessFadingmodel", "", "the name of access intf fadingmodel")
	kubenetCmd.Flags().String("accessLinkmodel", "", "the name of access intf linkmodel")
	kubenetCmd.Flags().String("accessTopologymodel", "", "the name of access intf topologymodel")

	// params for input and output
	kubenetCmd.Flags().StringP("template", "t", "./test.template.yaml", "Kubenet template YAML input path")
	kubenetCmd.Flags().StringP("input", "i", "./test.tle", "TLE input path")
	kubenetCmd.Flags().StringP("output", "o", "./test.yaml", "Kubenet YAML output path")
}
