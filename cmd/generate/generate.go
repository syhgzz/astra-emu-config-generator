package generate

import (
	"dslab.sjtu/kubenet-config-generator/cmd/generate/ground"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// GenerateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// GenerateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
