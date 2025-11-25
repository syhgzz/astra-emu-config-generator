package ground

import (
	"fmt"
	"os"

	"dslab.sjtu/kubenet-config-generator/internal/geo"
	"github.com/spf13/cobra"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Generate Kubenet-Compatible path for ground",
	Long:  "Generate Kubenet-Compatible path for ground.",
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, _ := cmd.Flags().GetString("file")
		content, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		output, err := geo.ConvertPathFromSAS(string(content), 0.0, 60, "yaml")
		if err != nil {
			fmt.Printf("Error converting path: %v\n", err)
			return
		}

		fmt.Println(output)
	},
}

func init() {
	GroundCmd.AddCommand(pathCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pathCmd.PersistentFlags().String("foo", "", "A help for foo")

	pathCmd.PersistentFlags().StringP("file", "f", "", "Input file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pathCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
