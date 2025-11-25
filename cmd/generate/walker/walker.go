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
	// DO NOT SET DEFAULTS HERE
	WalkerCmd.PersistentFlags().StringP("config", "c", "", "Path to YAML config file")
	WalkerCmd.PersistentFlags().IntP("shell-idx", "s", 1, "the index of shell")
	_ = WalkerCmd.MarkFlagRequired("shell-idx")
	WalkerCmd.PersistentFlags().IntP("num", "N", 0, "Total number of satellites, which must be divisible by P")
	WalkerCmd.PersistentFlags().IntP("plane", "P", 0, "The num of orbital planes")
	WalkerCmd.PersistentFlags().IntP("needed-plane", "E", 0, "The num of plane used")
	WalkerCmd.PersistentFlags().Float64P("phase", "F", 0, "Phasing between satellites in adjacent planes, range from 0 to plane-1")
}
