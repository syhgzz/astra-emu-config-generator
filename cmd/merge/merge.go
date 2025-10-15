package merge

import (
	"fmt"
	"os"
	"path/filepath"

	"dslab.sjtu/kubenet-config-generator/pkg/kubenet"
	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// tleCmd represents the tle command
var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge yaml files",
	Long:  `Merge yaml files.`,
	Run: func(cmd *cobra.Command, args []string) {
		intputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")

		// Check if input path exists
		if _, err := os.Stat(intputPath); os.IsNotExist(err) {
			logrus.Fatalf("input path does not exist: %v", intputPath)
			return
		}

		entries, err := os.ReadDir(intputPath)
		if err != nil {
			logrus.Fatalf("failed to read directory: %v", err)
			return
		}

		networkExperiments := []*kubenet_config.NetworkExperiment{}
		for _, entry := range entries {
			filename := filepath.Join(intputPath, entry.Name())
			fmt.Println("merge file:", filename)

			entryNetworkExperiment, err := kubenet.ReadConfig(filename)
			if err != nil {
				logrus.Fatalf("failed to unmarshal networkExperiment: %v", err)
			}
			networkExperiments = append(networkExperiments, entryNetworkExperiment)

		}

		networkExperiment := kubenet.MergeNetworkExperiment(networkExperiments)

		err = kubenet.WriteConfig(outputPath, networkExperiment)
		if err != nil {
			logrus.Fatalf("failed to write merged networkExperiment: %v", err)
		}
	},
}

func init() {

	MergeCmd.Flags().StringP("input", "i", "./data", "directory of yaml files to be merged")
	MergeCmd.Flags().StringP("output", "o", "./testAll.yaml", "merged file output path")
}
