package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubenet-config-generator",
	Short: "A tool to generate kubenet config",
	Long: `A tool to generate kubenet config

Examples:

# generate 720 satellites with 80 planes, the default output is test.tle
go run . tle -N 720 -P 80

# generate kubenet config yaml file from tle file
go run . generate -N 720 -P 80 -i test.tle -o test.yaml
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
