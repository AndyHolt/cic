/*
Copyright © 2024 Andy Holt <andrew.holt@hotmail.co.uk>
*/
package cmd

import (
	"github.com/AndyHolt/cic/imgproc"

	"github.com/spf13/cobra"
)

var Clusters int

// kmeansCmd represents the kmeans command
var kmeansCmd = &cobra.Command{
	Use:   "kmeans",
	Short: "Perform k-means clustering on an image file",
	Long: `Perform k-means clustering (vector quantisation) on an image file.

This reduces information in the image by a means of compression and
dimensionality reduction, and may lead to distinct areas of the image being
better identified by edge detection algorithms.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cic.RunKMeansImage(args[0], OutputFileName, Clusters)
	},
}

func init() {
	rootCmd.AddCommand(kmeansCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kmeansCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kmeansCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	kmeansCmd.Flags().IntVarP(&Clusters, "clusters", "k", 4,
		"Number of clusters for k-means")
}
