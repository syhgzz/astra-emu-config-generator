package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"dslab.sjtu/kubenet-config-generator/api/pkg/kubenet_config"
	"dslab.sjtu/kubenet-config-generator/pkg/kubenet"
	"dslab.sjtu/kubenet-config-generator/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// tleCmd represents the tle command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge yaml files",
	Long:  `Merge yaml files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start merge yaml files")
		intputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")

		entries, err := os.ReadDir(intputPath)
		if err != nil {
			logrus.Fatalf("failed to read directory: %v", err)
			return
		}

		networkExperiment := &kubenet_config.NetworkExperiment{}
		for _, entry := range entries {
			filename := filepath.Join(intputPath, entry.Name())
			fmt.Println("merge file:", filename)
			yamlData, err := os.ReadFile(filename)
			if err != nil {
				logrus.Fatalf("failed to read networkExperiment: %v", err)
			}
			entryNetworkExperiment, err := util.YAMLToStruct(yamlData)
			if err != nil {
				logrus.Fatalf("failed to unmarshal networkExperiment: %v", err)
			}
			kubenet.MergeNetworkExperiment(networkExperiment, entryNetworkExperiment)
		}

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
	rootCmd.AddCommand(mergeCmd)

	mergeCmd.Flags().StringP("input", "i", "./data", "directory of yaml files to be merged")
	mergeCmd.Flags().StringP("output", "o", "./testAll.yaml", "merged file output path")
}
