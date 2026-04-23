package simple

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/simple"
	"github.com/spf13/cobra"
)

var SimpleCmd = &cobra.Command{
	Use:   "simple",
	Short: "Generate Simple Config",
	Long:  `Generate Simple Config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("✏️  Start generating Kubenet config")

		configFile, _ := cmd.Flags().GetString("config")
		outputPath, _ := cmd.Flags().GetString("output")

		c, err := simple.LoadSimpleConfigFromYAML(configFile)
		if err != nil {
			return fmt.Errorf("failed to load kubenet generator config from config file %s: %w", configFile, err)
		}

		generator := simple.NewSimpleGenerator(c, outputPath)

		if err := generator.Generate(); err != nil {
			fmt.Println("❌ generate kubenet config failed")
			return err
		}

		fmt.Println("✅ Kubenet config generated:", outputPath)
		return nil
	},
}

func init() {
	// DO NOT SET DEFAULTS HERE
	// defaults are set in NewDefaultGeneratorConfig
	// GenerateCmd.AddCommand(SimpleCmd)

	// params for input and output
	SimpleCmd.PersistentFlags().StringP("config", "c", "", "Path to YAML config file")
	SimpleCmd.PersistentFlags().StringP("output", "o", "", "Kubenet YAML output path")
	_ = SimpleCmd.MarkFlagRequired("output")
}
