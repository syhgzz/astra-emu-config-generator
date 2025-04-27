package cmd

import (
	"fmt"
	"os"
	"strconv"

	"dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"
	"dslab.sjtu/kubenet-config-generator/pkg/kubenet"
	"dslab.sjtu/kubenet-config-generator/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// kubenetCmd represents the generate command
var kubenetCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Kubenet Config YAML",
	Long:  `Generate Kubenet Config YAML.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start generating Kubenet Config YAML")
		num, _ := cmd.Flags().GetInt("num")
		plane, _ := cmd.Flags().GetInt("plane")
		inputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")
		bandwidth, _ := cmd.Flags().GetString("bandwidth")
		singleHopLatency, _ := cmd.Flags().GetString("singleHopLatency")

		satNumPerTrack := num / plane

		constellation, err := kubenet.NewConstellation(inputPath)
		if err != nil {
			logrus.Fatalf("failed to initialize constellation: %v", err)
		}

		linkModels := []kubenet_config.LinkModel{}
		topologyModels := []kubenet_config.TopologyModel{}
		networkNodes := []kubenet_config.NetworkNode{}
		routingAreas := []kubenet_config.RoutingArea{}
		constellations := []kubenet_config.Constellation{}
		extra := kubenet_config.Extra{}

		walkerInterSatTopoIntfParam := kubenet_config.WalkerInterSatTopoIntfParam{
			Type: "WalkerInterSatTopoIntf",
		}
		topologyParam := kubenet_config.IntfTopologyParam{
			WalkerInterSatTopoIntfParam: &walkerInterSatTopoIntfParam,
		}
		satelliteLaserLinkIntfParam := kubenet_config.SatelliteLaserLinkIntfParam{
			Type: "SatelliteLaserLinkIntf",
		}
		linkParam := kubenet_config.IntfLinkParam{
			SatelliteLaserLinkIntfParam: &satelliteLaserLinkIntfParam,
		}

		containerParam := kubenet_config.ContainerParam{
			Type: "Container",
		}
		networkNodeParam := kubenet_config.NetworkNodeParam{
			ContainerParam: &containerParam,
		}

		// Initialize LinkModel
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

		// Initialize TopologyModel
		walkerInterSatTopoParam := kubenet_config.WalkerInterSatTopoParam{
			Type: "WalkerInterSatTopoParam",
		}
		topologyModelParam := kubenet_config.TopologyModelParam{
			WalkerInterSatTopoParam: &walkerInterSatTopoParam,
		}
		topologyModel := kubenet_config.TopologyModel{
			Name:  "topologymodel1",
			Param: topologyModelParam,
		}
		topologyModels = append(topologyModels, topologyModel)

		// Initialize RoutingArea
		OSPFV3Param := kubenet_config.OSPFV3Param{
			Protocol:     "OSPFV3",
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

		// Initialize Constellations and NetworkNodes
		satellites := []kubenet_config.OrbitSatellite{}
		planes := []kubenet_config.OrbitPlane{}
		planeIdx := 0
		innerIdx := 0
		for idx := range num {
			satellites = append(
				satellites,
				kubenet_config.OrbitSatellite{
					Idx:       int32(innerIdx),
					Satellite: constellation[idx].Name,
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
				TleLine1: constellation[idx].Satellite.Line1,
				TleLine2: constellation[idx].Satellite.Line2,
			}
			intf := []kubenet_config.Intf{}
			for i := 0; i < 4; i++ {
				intf = append(
					intf,
					kubenet_config.Intf{
						Name:          "intf" + strconv.Itoa(i),
						TopologyModel: "topologymodel1",
						LinkModel:     "linkmodel1",
						TopologyParam: topologyParam,
						LinkParam:     linkParam,
						RoutingArea:   &routingArea.Id,
					},
				)
			}
			mobilityModel := kubenet_config.NetworkNodeMobilityModel{
				SatelliteMobilityModel: &satelliteMobilityModel,
			}
			visualOptions := kubenet_config.VisualOption{
				ModelId:  1,
				NodeType: 1,
			}
			networkNodes = append(
				networkNodes,
				kubenet_config.NetworkNode{
					Name:          constellation[idx].Name,
					Type:          "Container",
					VisualOptions: &visualOptions,
					Intfs:         intf,
					MobilityModel: mobilityModel,
					Param:         networkNodeParam,
				},
			)
		}
		walkerDeltaParam := kubenet_config.WalkerDeltaParam{
			Type:   "WalkerDelta",
			Planes: planes,
		}
		spheres := []kubenet_config.ConstellationSpheresInner{}
		spheres = append(
			spheres,
			kubenet_config.ConstellationSpheresInner{
				WalkerDeltaParam: &walkerDeltaParam,
			},
		)
		constellations = append(
			constellations,
			kubenet_config.Constellation{
				Name:    "WalkerDelta",
				Spheres: spheres,
			},
		)

		// Construct NetworkExperiment and convert it to yaml file
		networkExperiment := kubenet_config.NetworkExperiment{
			LinkModels:     linkModels,
			TopologyModels: topologyModels,
			NetworkNodes:   networkNodes,
			RoutingAreas:   routingAreas,
			Constellations: constellations,
			Extra:          &extra,
		}

		// bytes, err := networkExperiment.MarshalJSON()
		bytes, err := util.StructToYAML(&networkExperiment)
		if err != nil {
			logrus.Fatalf("failed to marshal networkExperiment: %v", err)
		}
		err = os.WriteFile(outputPath, bytes, 0644)
		if err != nil {
			logrus.Fatalf("failed to write networkExperiment to file: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(kubenetCmd)

	// params for walker-delta constellation model
	kubenetCmd.Flags().IntP("num", "N", 300, "Total number of satellites, which must be divisible by P")
	kubenetCmd.Flags().IntP("plane", "P", 10, "The num of orbital planes")

	// params for network experiment
	kubenetCmd.Flags().StringP("bandwidth", "b", "1Gibps", "linkmodel Bandwidth")
	kubenetCmd.Flags().StringP("singleHopLatency", "s", "5ms", "linkmodel SingleHopLatency")

	// params for input and output
	kubenetCmd.Flags().StringP("input", "i", "./test.tle", "TLE input path")
	kubenetCmd.Flags().StringP("output", "o", "./test.yaml", "Kubenet YAML output path")
}
