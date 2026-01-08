package walker

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/cmd/util"
	"dslab.sjtu/kubenet-config-generator/pkg/walker/kubenet"
	"github.com/spf13/cobra"
)

// kubenetCmd represents the generate command
var kubenetCmd = &cobra.Command{
	Use:   "kubenet",
	Short: "Generate Kubenet Config",
	Long:  `Generate Kubenet Config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("✏️  Start generating Kubenet config")

		configFile, _ := cmd.Flags().GetString("config")
		shellIdx, _ := cmd.Flags().GetInt("shell-idx")
		outputPath, _ := cmd.Flags().GetString("output")
		inputPath, _ := cmd.Flags().GetString("input")
		templatePath, _ := cmd.Flags().GetString("template")

		c := kubenet.NewDefaultGeneratorConfig()

		if configFile != "" {
			// 优先从 YAML 读取
			err := kubenet.LoadGeneratorConfigFromYAML(configFile, shellIdx, &c)
			if err != nil {
				return fmt.Errorf("failed to load kubenet generator config from config file %s: %w", configFile, err)
			}
			fmt.Println("📖 Loaded kubenet generator config from YAML:", configFile)
		} else {
			util.SetIfChanged[int](cmd, "shell-idx", func(v int) { c.ShellIdx = v })
			util.SetIfChanged[int](cmd, "num", func(v int) { c.SatelliteNum = v })
			util.SetIfChanged[int](cmd, "plane", func(v int) { c.PlaneNum = v })
			util.SetIfChanged[int](cmd, "needed-plane", func(v int) { c.NeededPlane = v })
			util.SetIfChanged[string](cmd, "node-name-prefix", func(v string) { c.NodeNamePrefix = v })
			util.SetIfChanged[int](cmd, "access-intf-num", func(v int) { c.AccessIntfNum = v })
			util.SetIfChanged[int32](cmd, "offset", func(v int32) { c.Offset = v })
			util.SetIfChanged[string](cmd, "image", func(v string) { c.Image = v })
			util.SetIfChanged[int32](cmd, "routing-area", func(v int32) { c.RoutingArea = v })
			util.SetIfChanged[bool](cmd, "fake", func(v bool) { c.Fake = v })
			util.SetIfChanged[string](cmd, "gen-flow", func(v string) { c.GenFlow = v })
			util.SetIfChanged[bool](cmd, "cross-sphere", func(v bool) { c.CrossSphere = v })
			util.SetIfChanged[bool](cmd, "ground-access", func(v bool) { c.GroundAccess = v })
			util.SetIfChanged[float64](cmd, "phase", func(v float64) { c.Phase = v })
		}

		generator := kubenet.NewKubenetGenerator(
			kubenet.WithConfig(c),
			kubenet.WithInputPath(inputPath),
			kubenet.WithTemplatePath(templatePath),
			kubenet.WithOutputPath(outputPath),
		)

		if err := generator.Generate(); err != nil {
			fmt.Println("❌ generate kubenet config failed")
			return err
		}

		fmt.Println("✅ Kubenet config generated:", outputPath)
		return nil
	},
}

func init() {
	WalkerCmd.AddCommand(kubenetCmd)

	// DO NOT SET DEFAULTS HERE
	// defaults are set in NewDefaultGeneratorConfig

	// params for constellation
	kubenetCmd.Flags().Int32("offset", 0, "Offset")
	kubenetCmd.Flags().Bool("cross-sphere", false, "cross-sphere")

	// params for network experiment
	kubenetCmd.Flags().String("node-name-prefix", "", "the prefix of network node name")
	kubenetCmd.Flags().String("image", "", "application image")
	kubenetCmd.Flags().Int32P("routing-area", "r", 0, "routing area Id")
	kubenetCmd.Flags().Bool("fake", false, "Fake")
	kubenetCmd.Flags().String("gen-flow", "", "how to generate flow tasks")

	// params for ground access intfs
	kubenetCmd.Flags().Bool("ground-access", false, "enable ground access")
	kubenetCmd.Flags().Int("access-intf-num", 0, "the number of access intfs")

	// params for input and output
	kubenetCmd.Flags().StringP("template", "t", "", "Kubenet template YAML input path")
	kubenetCmd.Flags().StringP("input", "i", "", "TLE input path")
	kubenetCmd.Flags().StringP("output", "o", "", "Kubenet YAML output path")
	_ = kubenetCmd.MarkFlagRequired("template")
	_ = kubenetCmd.MarkFlagRequired("input")
	_ = kubenetCmd.MarkFlagRequired("output")
}
