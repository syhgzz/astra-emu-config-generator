package cmd

import (
	"fmt"
	"os"
	"strconv"

	"dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"
	"dslab.sjtu/kubenet-config-generator/pkg/kubenet"
	"dslab.sjtu/kubenet-config-generator/pkg/tle"
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
		nodeNamePre, _ := cmd.Flags().GetString("nodeNamePre")
		template, _ := cmd.Flags().GetString("template")
		inputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")

		satNumPerTrack := num / plane

		yamlData, err := os.ReadFile(template)
		if err != nil {
			logrus.Fatalf("failed to read networkExperiment: %v", err)
		}
		networkExperiment, err := util.YAMLToStruct(yamlData)
		if err != nil {
			logrus.Fatalf("failed to unmarshal networkExperiment: %v", err)
		}
		networkNodes := []kubenet_config.NetworkNode{}
		constellations := []kubenet_config.Constellation{}

		phyConstellation, err := tle.NewConstellation(inputPath)
		if err != nil {
			logrus.Fatalf("failed to initialize constellation: %v", err)
		}

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

		routingArea := networkExperiment.RoutingAreas[0]

		visualOptions := kubenet_config.VisualOption{
			ModelId:  1,
			NodeType: 1,
		}

		// Initialize Constellations and NetworkNodes
		satellites := []kubenet_config.OrbitSatellite{}
		planes := []kubenet_config.OrbitPlane{}
		planeIdx := 0
		innerIdx := 0
		for idx := range num {
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
			for i := 0; i < 4; i++ {
				intfs = append(
					intfs,
					kubenet_config.Intf{
						Name:          "intf" + strconv.Itoa(i),
						TopologyModel: "walkerinter",
						LinkModel:     "satlaser",
						TopologyParam: topologyParam,
						LinkParam:     linkParam,
						RoutingArea:   &routingArea.Id,
					},
				)
			}
			intfs = append(
				intfs,
				kubenet.NewSatConeIntf(routingArea),
			)
			mobilityModel := kubenet_config.NetworkNodeMobilityModel{
				SatelliteMobilityModel: &satelliteMobilityModel,
			}

			networkNodes = append(
				networkNodes,
				kubenet_config.NetworkNode{
					Name:          nodeName,
					Type:          "Container",
					VisualOptions: &visualOptions,
					Intfs:         intfs,
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

		networkNodes = append(networkNodes, kubenet.NewGroundConeNetworkNode(routingArea))

		// Construct NetworkExperiment and convert it to yaml file
		networkExperiment.Constellations = append(networkExperiment.Constellations, constellations...)
		networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, networkNodes...)

		// bytes, err := networkExperiment.MarshalJSON()
		bytes, err := util.StructToYAML(networkExperiment)
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
	kubenetCmd.Flags().StringP("nodeNamePre", "n", "", "the prefix of network node name")
	kubenetCmd.Flags().StringP("bandwidth", "b", "1Gibps", "linkmodel Bandwidth")
	kubenetCmd.Flags().StringP("singleHopLatency", "s", "5ms", "linkmodel SingleHopLatency")

	// params for input and output
	kubenetCmd.Flags().StringP("template", "t", "./test.template.yaml", "Kubenet template YAML input path")
	kubenetCmd.Flags().StringP("input", "i", "./test.tle", "TLE input path")
	kubenetCmd.Flags().StringP("output", "o", "./test.yaml", "Kubenet YAML output path")
}
