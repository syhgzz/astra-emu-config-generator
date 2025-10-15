package cmd

import (
	"os"

	"dslab.sjtu/kubenet-config-generator/cmd/generate"
	"dslab.sjtu/kubenet-config-generator/cmd/merge"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubenet-config-generator",
	Short: "A tool to operate on kubenet config",
	Long:  `A tool to operate on kubenet config`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(generate.GenerateCmd)
	RootCmd.AddCommand(merge.MergeCmd)
}
