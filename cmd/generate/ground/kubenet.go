package ground

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet_config_util"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// kubenetCmd represents the ground command
var kubenetCmd = &cobra.Command{
	Use:   "kubenet",
	Short: "Generate Kubenet config for ground",
	Long:  `Generate Kubenet config for ground.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start generating Ground cluster YAML")
		num, _ := cmd.Flags().GetInt("num")
		name, _ := cmd.Flags().GetString("name")
		longitude, _ := cmd.Flags().GetFloat64("longitude")
		latitude, _ := cmd.Flags().GetFloat64("latitude")
		template, _ := cmd.Flags().GetString("template")
		outputPath, _ := cmd.Flags().GetString("output")
		accessFadingmodel, _ := cmd.Flags().GetString("accessFadingmodel")
		accessLinkmodel, _ := cmd.Flags().GetString("accessLinkmodel")
		accessTopologymodel, _ := cmd.Flags().GetString("accessTopologymodel")
		image, _ := cmd.Flags().GetString("image")
		routingAreaId, _ := cmd.Flags().GetInt32("routingArea")

		var routingAreaPtr *int32
		if routingAreaId >= 0 {
			routingAreaPtr = &routingAreaId
		}

		networkExperiment, err := kubenet_config_util.ReadConfig(template)
		if err != nil {
			logrus.Fatalf("failed to read networkExperiment: %v", err)
		}

		networkNodes := []kubenet_config.NetworkNode{}

		coneTopoIntfParam := kubenet_config.SatelliteConeTopoIntfParam{
			Type: "SatelliteConeTopoIntf",
			Role: 0,
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
			ModelId:  "4",
			NodeType: 2,
		}

		earthCenteredPointingOrientationModel := kubenet_config.EarthCenteredPointingOrientationModel{
			Type: "EarthCenteredPointing",
		}

		orientationModel := kubenet_config.NetworkNodeOrientationModel{
			EarthCenteredPointingOrientationModel: &earthCenteredPointingOrientationModel,
		}

		for idx := range num {
			nodeName := strings.ToLower(name) + strconv.Itoa(idx)

			intfs := []kubenet_config.Intf{}
			intfs = append(
				intfs,
				kubenet_config.Intf{
					Name:          "acce1",
					TopologyModel: accessTopologymodel,
					LinkModel:     accessLinkmodel,
					TopologyParam: coneTopologyParam,
					LinkParam:     accessIntfLinkParam,
					RoutingArea:   routingAreaPtr,
				},
			)

			// TODO
			staticMobilityModel := kubenet_config.StaticMobilityModel{
				Type:      "Static",
				Longitude: longitude + rand.Float64()*2.0,
				Latitude:  latitude + rand.Float64()*2.0,
				Altitude:  50,
			}

			mobilityModel := kubenet_config.NetworkNodeMobilityModel{
				StaticMobilityModel: &staticMobilityModel,
			}

			networkNodes = append(
				networkNodes,
				kubenet_config.NetworkNode{
					Name:             nodeName,
					Type:             "Container",
					VisualOptions:    &visualOptions,
					Intfs:            intfs,
					MobilityModel:    mobilityModel,
					OrientationModel: orientationModel,
					Param:            &networkNodeParam,
				},
			)
		}

		networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, networkNodes...)

		err = kubenet_config_util.WriteConfig(outputPath, networkExperiment)
		if err != nil {
			logrus.Fatalf("failed to write networkExperiment: %v", err)
		}
	},
}

func init() {
	GroundCmd.AddCommand(kubenetCmd)

	// params for ground cluster
	kubenetCmd.Flags().StringP("name", "n", "", "the name of ground cluster")
	kubenetCmd.Flags().IntP("num", "N", 20, "Total number of ground cluster nodes")
	kubenetCmd.Flags().Float64("longitude", 0.0, "longitude")
	kubenetCmd.Flags().Float64("latitude", 0.0, "latitude")

	// params for network experiment
	kubenetCmd.Flags().String("image", "", "application image")
	kubenetCmd.Flags().Int32P("routingArea", "r", -1, "routing area Id")

	// params for ground access intfs
	kubenetCmd.Flags().String("accessFadingmodel", "fading_model1", "the name of access intf fadingmodel")
	kubenetCmd.Flags().String("accessLinkmodel", "Satellite", "the name of access intf linkmodel")
	kubenetCmd.Flags().String("accessTopologymodel", "cone1", "the name of access intf topologymodel")

	// params for input and output
	kubenetCmd.Flags().StringP("template", "t", "./test.template.yaml", "Kubenet template YAML input path")
	kubenetCmd.Flags().StringP("output", "o", "./ground.yaml", "Kubenet YAML output path")
}
