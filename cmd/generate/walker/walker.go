package walker

import (
	"github.com/spf13/cobra"
)

// WalkerCmd represents the walker command
var WalkerCmd = &cobra.Command{
	Use:   "walker",
	Short: "Walker constellation functions",
	Long:  "Walker constellation functions.",
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walkerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// params for walker constellation model
	WalkerCmd.PersistentFlags().StringP("name", "n", "starlink", "the name of constellation")
	WalkerCmd.PersistentFlags().IntP("num", "N", 300, "Total number of satellites, which must be divisible by P")
	WalkerCmd.PersistentFlags().IntP("plane", "P", 10, "The num of orbital planes")
	WalkerCmd.PersistentFlags().IntP("needed-plane", "E", 0, "The num of plane used")
	WalkerCmd.PersistentFlags().Float64P("phase", "F", 0, "Phasing between satellites in adjacent planes, range from 0 to plane-1")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walkerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
