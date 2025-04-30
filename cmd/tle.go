package cmd

import (
	"fmt"

	"dslab.sjtu/kubenet-config-generator/pkg/tle"
	"github.com/spf13/cobra"
)

// tleCmd represents the tle command
var tleCmd = &cobra.Command{
	Use:   "tle",
	Short: "Generate TLE",
	Long:  `Generate TLE from params.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start generating TLE")
		height, _ := cmd.Flags().GetFloat64("height")
		num, _ := cmd.Flags().GetInt("num")
		plane, _ := cmd.Flags().GetInt("plane")
		neededPlane, _ := cmd.Flags().GetInt("needed-plane")
		deltaInit, _ := cmd.Flags().GetFloat64("delta-init")
		phase, _ := cmd.Flags().GetFloat64("phase")
		inclination, _ := cmd.Flags().GetFloat64("inclination")
		outputPath, _ := cmd.Flags().GetString("output")

		if height <= 0 || num <= 0 || plane <= 0 {
			fmt.Println("height, num, plane must be positive")
			return
		}
		if num%plane != 0 {
			fmt.Println("num must be divisible by plane")
			return
		}
		if inclination < 0 || inclination > 90 {
			fmt.Println("inclination must be between 0 and 90")
			return
		}
		if phase < 0 || phase >= float64(plane) {
			fmt.Println("phase must be between 0 and plane")
			return
		}

		tleParams := tle.TLEParams{
			Height:      height,
			Num:         num,
			Plane:       plane,
			Phase:       phase,
			Inclination: inclination,

			DeltaInit:   deltaInit,
			NeededPlane: neededPlane,
		}
		items := tleParams.NewTLE()

		err := items.WriteFile(outputPath)
		if err != nil {
			fmt.Printf("Failed to write TLE file: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tleCmd)

	tleCmd.Flags().Float64P("height", "H", 500, "Orbit height(km), recommended 300km - 1500km")
	tleCmd.Flags().IntP("num", "N", 300, "Total number of satellites, which must be divisible by P")
	tleCmd.Flags().IntP("plane", "P", 10, "The num of orbital planes")
	tleCmd.Flags().IntP("needed-plane", "E", 0, "The num of plane used")
	tleCmd.Flags().IntP("delta-init", "D", 0, "The initial delta of the first plane")
	tleCmd.Flags().Float64P("phase", "F", 0, "Phasing between satellites in adjacent planes, range from 0 to plane-1")
	tleCmd.Flags().Float64P("inclination", "I", 50, "Orbital inclination")

	tleCmd.Flags().StringP("output", "o", "./test.tle", "TLE output path")
}
