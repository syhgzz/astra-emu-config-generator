package merge

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/merge"
	"github.com/spf13/cobra"
)

// tleCmd represents the tle command
var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge yaml files",
	Long:  `Merge yaml files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath, _ := cmd.Flags().GetString("input")
		outputPath, _ := cmd.Flags().GetString("output")

		fmt.Println("🔀 Start merging yaml files:", inputPath, "➡️ ", outputPath)

		merger := merge.NewConfigMerger(
			merge.WithInputPath(inputPath),
			merge.WithOutputPath(outputPath),
			merge.WithRecursive(false),
			merge.WithFilterExt(".yaml"),
			merge.WithIgnoreError(false),
		)

		if err := merger.Merge(); err != nil {
			return fmt.Errorf("merge config failed: %w", err)
		}
		return nil
	},
}

func init() {
	MergeCmd.Flags().StringP("input", "i", "./data", "directory of yaml files to be merged")
	_ = MergeCmd.MarkFlagRequired("input")
	MergeCmd.Flags().StringP("output", "o", "./testAll.yaml", "merged file output path")
	_ = MergeCmd.MarkFlagRequired("output")
}
