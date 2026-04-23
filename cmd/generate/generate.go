package generate

import (
	"dslab.sjtu/kubenet-config-generator/cmd/generate/ground"
	"dslab.sjtu/kubenet-config-generator/cmd/generate/simple"
	"dslab.sjtu/kubenet-config-generator/cmd/generate/walker"
	"github.com/spf13/cobra"
)

// GenerateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate configurations or data",
	Long:  `Generate configurations or data which used for construct experiments`,
}

func init() {
	GenerateCmd.AddCommand(ground.GroundCmd)
	GenerateCmd.AddCommand(walker.WalkerCmd)
	GenerateCmd.AddCommand(simple.SimpleCmd)
}
