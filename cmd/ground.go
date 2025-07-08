package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"dslab.sjtu/kubenet-config-generator/pkg/util"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// groundCmd represents the ground command
var groundCmd = &cobra.Command{
	Use:   "ground",
	Short: "Generate ground cluster",
	Long:  `Generate ground cluster from params.`,
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

		yamlData, err := os.ReadFile(template)
		if err != nil {
			logrus.Fatalf("failed to read networkExperiment: %v", err)
		}
		networkExperiment, err := util.YAMLToStruct(yamlData)
		if err != nil {
			logrus.Fatalf("failed to unmarshal networkExperiment: %v", err)
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

		ports := []int32{5000, 5001, 5002}
		applications := []kubenet_config.Application{}
		applications = append(
			applications,
			kubenet_config.Application{
				Image: "harbor.sail.se.sjtu.edu.cn/utils/test-suite:v0.2",
				Ports: ports,
			},
		)
		if image != "" {
			applications = append(
				applications,
				kubenet_config.Application{
					Image: image,
				},
			)
		}
		containerParam := kubenet_config.ContainerParam{
			Type:         "Container",
			Applications: applications,
		}
		networkNodeParam := kubenet_config.NetworkNodeParam{
			ContainerParam: &containerParam,
		}

		visualOptions := kubenet_config.VisualOption{
			ModelId:  "4",
			NodeType: 2,
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
					RoutingArea:   &routingAreaId,
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
					Name:          nodeName,
					Type:          "Container",
					VisualOptions: &visualOptions,
					Intfs:         intfs,
					MobilityModel: mobilityModel,
					Param:         &networkNodeParam,
				},
			)
		}

		networkExperiment.NetworkNodes = append(networkExperiment.NetworkNodes, networkNodes...)

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
	rootCmd.AddCommand(groundCmd)

	// params for ground cluster
	groundCmd.Flags().StringP("name", "n", "", "the name of ground cluster")
	groundCmd.Flags().IntP("num", "N", 20, "Total number of ground cluster nodes")
	groundCmd.Flags().Float64("longitude", 0.0, "longitude")
	groundCmd.Flags().Float64("latitude", 0.0, "latitude")

	// params for network experiment
	groundCmd.Flags().String("image", "", "application image")
	groundCmd.Flags().Int32P("routingArea", "r", 1, "routing area Id")

	// params for ground access intfs
	groundCmd.Flags().String("accessFadingmodel", "", "the name of access intf fadingmodel")
	groundCmd.Flags().String("accessLinkmodel", "", "the name of access intf linkmodel")
	groundCmd.Flags().String("accessTopologymodel", "", "the name of access intf topologymodel")

	// params for input and output
	groundCmd.Flags().StringP("template", "t", "./test.template.yaml", "Kubenet template YAML input path")
	groundCmd.Flags().StringP("output", "o", "./ground.yaml", "Kubenet YAML output path")
}
