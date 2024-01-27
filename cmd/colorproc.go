/*
Copyright © 2024 Andy Holt <andrew.holt@hotmail.co.uk
*/
package cmd

import (
	"github.com/AndyHolt/cic/imgproc"

	"github.com/spf13/cobra"
)

// colorprocCmd represents the colorproc command
var colorprocCmd = &cobra.Command{
	Use:   "colorproc",
	Short: "Testing processing image using colour information",
	Long: `Run Gaussian blur and edge detection using colour information

Standard usage of cic converts picture to grayscale, making use only of colour
intensity for further processing. This is an experimental feature to use colour
information for better edge detection between regions of similar colour
intensity, but different colour profile.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cic.RunColourImageProc(args[0], OutputFileName)
	},
}

func init() {
	rootCmd.AddCommand(colorprocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// colorprocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// colorprocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
