/*
Copyright © 2023 Andy Holt <andrew.holt@hotmail.co.uk>
*/
package cmd

import (
	"os"

	"github.com/AndyHolt/cic/imgproc"
	"github.com/spf13/cobra"
)

var OutputFileName string
var StdDev float64
var UpperThreshold int
var LowerThreshold int
var NonMaxSuppressionDistance int

// upper (u) and lower (l) for setting threshold values
// no-blur option

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cic [flags] filename",
	Short: "CIC is a Colouring-In Creator (or Colouring Image Creator)",
	Long: `CIC, the Colouring-In Creator, turns images into colouring sheets.

CIC uses image processing and edge detection techniques to turn any image file
into a colouring sheet.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cic.ConvertImageToColouring(
			args[0],
			OutputFileName,
			StdDev,
			UpperThreshold,
			LowerThreshold,
			NonMaxSuppressionDistance,
		)
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cic.yaml)")
	rootCmd.PersistentFlags().StringVarP(&OutputFileName, "output", "o",
		"edited.jpg", "File name of output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Float64VarP(&StdDev, "stddev", "s", 1.0,
		"Std dev for Gaussian blur")
	rootCmd.Flags().IntVarP(&UpperThreshold, "upper", "u", 100,
		"Upper threshold for edge suppression")
	rootCmd.Flags().IntVarP(&LowerThreshold, "lower", "l", 10,
		"Lower threshold for edge suppression")
	rootCmd.Flags().IntVarP(&NonMaxSuppressionDistance, "distance", "d", 1,
		"Interval for non-maximum suppression in pixels")
}
