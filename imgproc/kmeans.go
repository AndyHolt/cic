package cic

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"

	_ "image/png"
)

type Pixel struct {
	X int
	Y int
}

type Mean struct {
	R uint8
	G uint8
	B uint8
}

type KMeansClusters struct {
	K        int
	Means    []Mean
	Clusters [][]Pixel
	CostVal  float64
}

func InitKMeans(k int) *KMeansClusters {
	var kmc KMeansClusters
	kmc.K = k
	kmc.Means = make([]Mean, k, k)
	kmc.Clusters = make([][]Pixel, k)
	kmc.CostVal = 0.0
	return &kmc
}

func (kmc *KMeansClusters) RandomiseMeans() {
	for i := 0; i < kmc.K; i++ {
		kmc.Means[i].R = uint8(rand.Intn(256))
		kmc.Means[i].G = uint8(rand.Intn(256))
		kmc.Means[i].B = uint8(rand.Intn(256))
	}
}

func (kmc *KMeansClusters) DistToMean(meanIdx int, px color.RGBA) float64 {
	return math.Sqrt(
		math.Pow(float64(px.R-kmc.Means[meanIdx].R), 2) +
			math.Pow(float64(px.G-kmc.Means[meanIdx].G), 2) +
			math.Pow(float64(px.B-kmc.Means[meanIdx].B), 2),
	)
}

func (kmc *KMeansClusters) AssignClusters(img *image.RGBA) {
	for i := 0; i < kmc.K; i++ {
		kmc.Clusters[i] = []Pixel{}
	}
	kmc.CostVal = 0.0

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			bestMean := -1
			bestMeanDist := math.MaxFloat64

			pxval := img.RGBAAt(x, y)

			for mIdx := 0; mIdx < kmc.K; mIdx++ {
				dist := kmc.DistToMean(mIdx, pxval)
				if dist < bestMeanDist {
					bestMean = mIdx
					bestMeanDist = dist
				}
			}

			kmc.Clusters[bestMean] = append(kmc.Clusters[bestMean], Pixel{x, y})
			kmc.CostVal += math.Pow(bestMeanDist, 2)
		}
	}
}

func (kmc *KMeansClusters) CalculateMeans(img *image.RGBA) {
	for i := 0; i < kmc.K; i++ {
		var runSum = [3]float64{0.0, 0.0, 0.0}
		norm := float64(len(kmc.Clusters[i]))
		for _, px := range kmc.Clusters[i] {
			pxval := img.RGBAAt(px.X, px.Y)
			runSum[0] += float64(pxval.R)
			runSum[1] += float64(pxval.G)
			runSum[2] += float64(pxval.B)
		}
		kmc.Means[i].R = uint8(runSum[0] / norm)
		kmc.Means[i].G = uint8(runSum[1] / norm)
		kmc.Means[i].B = uint8(runSum[2] / norm)
	}
}

func (kmc *KMeansClusters) AssignClusterMeanValues(img *image.RGBA) *image.RGBA {
	for meanIdx := 0; meanIdx < kmc.K; meanIdx++ {
		pxval := color.RGBA{
			kmc.Means[meanIdx].R,
			kmc.Means[meanIdx].G,
			kmc.Means[meanIdx].B,
			255,
		}
		for _, px := range kmc.Clusters[meanIdx] {
			img.SetRGBA(px.X, px.Y, pxval)
		}
	}
	return img
}

func KMeansImage(img *image.RGBA, k int) *image.RGBA {
	fmt.Printf("Running K-means with %v means", k)

	kmc := InitKMeans(k)

	// Assign initial random cluster means
	kmc.RandomiseMeans()

	// Assign pixels to clusters
	kmc.AssignClusters(img)

	fmt.Printf("Initial (random) setting of means gives cost: %v\n", kmc.CostVal)

	lastCostVal := kmc.CostVal

	// Main loop: iterate mean evaluation and cluster reassignment
	for i := 1; i <= 10; i++ {
		fmt.Printf("Beginning iteration: %v\nCalculating updated means...", i)
		kmc.CalculateMeans(img)
		fmt.Printf(" Done\n")
		fmt.Printf("Assigning pixels to clusters...")
		kmc.AssignClusters(img)
		fmt.Printf(" Done\n")
		fmt.Printf("On iteration %v, cost value is: %v\n", i, kmc.CostVal)
		fmt.Printf("Cost value improvement of: %v\n\n", kmc.CostVal-lastCostVal)
		lastCostVal = kmc.CostVal
	}

	fmt.Printf("Calculating final mean values for pixel assignment...")
	kmc.CalculateMeans(img)
	fmt.Printf(" Done\n")

	// Assign pixels to mean values and return modified image
	fmt.Printf("Setting pixel values based on cluster mean...")
	dstImg := kmc.AssignClusterMeanValues(img)
	fmt.Printf(" Done\n")

	return dstImg
}

func RunKMeansImage(filename string, outputFilename string, k int) {
	fmt.Print("Reading in file...")

	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	fmt.Print(" Done\n")
	fmt.Print("Decoding file to Image...")

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(" Done\n")

	fmt.Print("Convert to RBGA format...")
	rgba := imageToRGBA(img)
	fmt.Print(" Done\n")

	rgba = KMeansImage(rgba, k)

	fmt.Printf("Saving to output file: \"%v\"", outputFilename)
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	var imgOptions jpeg.Options
	imgOptions.Quality = 100

	jpeg.Encode(outputFile, rgba, &imgOptions)

	fmt.Print(" Done\n")
}
