package merge

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/merge"
	"github.com/spf13/cobra"
)

var inputPaths []string
var outputPath *string

// tleCmd represents the tle command
var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge kubenet config files (json or yaml)",
	Long: `
Merge kubenet config files (json or yaml) into one file (json or yaml, depending on the extension of the output file).
Example usages:
- $(CLI) merge -i config1.yaml,config2.yaml -o merged.yaml
- $(CLI) merge -i config1.yaml -i config2.yaml -o merged.yaml
- $(CLI) merge *.yaml -o merged.json
- $(CLI) merge -i *.yaml -o merged.yaml
- $(CLI) merge -i config1.yaml -o merged.yaml
- $(CLI) merge config1.yaml -o merged.yaml
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputPath == nil || *outputPath == "" {
			return fmt.Errorf("no output file specified")
		}

		// 1. 合并 -i 和 args
		allInputs := append(inputPaths, args...)

		if len(allInputs) == 0 {
			return fmt.Errorf("no input files specified")
		}

		// 2. 去重（避免重复）
		finalInputs := dedup(allInputs)

		fmt.Println("🔀 Start merging kubenet config files:", finalInputs, "➡️ ", outputPath)

		merger := merge.NewConfigMerger(
			merge.WithInputPaths(finalInputs),
			merge.WithOutputPath(*outputPath),
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
	MergeCmd.Flags().StringSliceVarP(
		&inputPaths,
		"input", "i",
		nil,
		"input yaml files (support multiple or comma-separated)",
	)
	outputPath = MergeCmd.Flags().StringP(
		"output", "o",
		"",
		"output file, must end with .yaml or .json",
	)
}

func dedup(inputs []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, in := range inputs {
		if _, ok := seen[in]; !ok {
			seen[in] = struct{}{}
			result = append(result, in)
		}
	}
	return result
}
