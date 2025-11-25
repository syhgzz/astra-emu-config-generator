package walker

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/cmd/util"
	"dslab.sjtu/kubenet-config-generator/pkg/walker/tle"
	"github.com/spf13/cobra"
)

// tleCmd represents the tle command
var tleCmd = &cobra.Command{
	Use:   "tle",
	Short: "Generate TLE",
	Long:  `Generate TLE from params.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Start generating TLE")

		configFile, _ := cmd.Flags().GetString("config")
		shellIdx, _ := cmd.Flags().GetInt("shell-idx")
		outputPath, _ := cmd.Flags().GetString("output")

		p := tle.NewDefaultTLEParams()

		if configFile != "" {
			// 优先从 YAML 读取
			err := tle.LoadTLEParamsFromYAML(configFile, shellIdx, &p)
			if err != nil {
				return fmt.Errorf("failed to load tle params from config file %s: %w", configFile, err)
			}
			fmt.Println("📖 Loaded tle params from YAML:", configFile)
		} else {
			// 从命令行读取参数
			util.SetIfChanged[int](cmd, "shell-idx", func(v int) { p.ConstellationIdx = v })
			util.SetIfChanged[float64](cmd, "height", func(v float64) { p.Height = v })
			util.SetIfChanged[int](cmd, "num", func(v int) { p.Num = v })
			util.SetIfChanged[int](cmd, "plane", func(v int) { p.Plane = v })
			util.SetIfChanged[int](cmd, "needed-plane", func(v int) { p.NeededPlane = v })
			util.SetIfChanged[float64](cmd, "delta-init", func(v float64) { p.DeltaInit = v })
			util.SetIfChanged[float64](cmd, "phase", func(v float64) { p.Phase = v })
			util.SetIfChanged[float64](cmd, "inclination", func(v float64) { p.Inclination = v })
			util.SetIfChanged[string](cmd, "name-rule", func(v string) { p.NameRule = v })
		}

		generator := tle.NewTLEGenerator(
			tle.WithParams(p),
			tle.WithOutput(outputPath),
		)

		if err := generator.Generate(); err != nil {
			fmt.Println("❌ generate tle failed")
			return err
		}

		fmt.Println("✅ TLE generated:", outputPath)
		return nil
	},
}

func init() {
	WalkerCmd.AddCommand(tleCmd)

	// DO NOT SET DEFAULTS HERE
	// defaults are set in NewDefaultTLEParams
	tleCmd.Flags().String("name-rule", "", "the name rule of constellation")
	tleCmd.Flags().Float64P("height", "H", 0, "Orbit height(km), recommended 300km - 1500km")
	tleCmd.Flags().Float64P("delta-init", "D", 0, "The initial delta of the first plane")
	tleCmd.Flags().Float64P("inclination", "I", 0, "Orbital inclination")
	tleCmd.Flags().StringP("output", "o", "", "TLE output path")
	_ = tleCmd.MarkFlagRequired("output")
}
