// Next steps:
// 1. Tidy up code.
// 2. Store size in ImageGradients datastructure so don't need to compute
//    repeatedly
// 3. Make parameters selectable (thresholds, Guassian blur size & sigma,
//    whether or not to use Gaussian blur),
// 4. Make CLI interface
// 5. Make a GUI/web interface to allow fast experimentation with parameter
//    tuning

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"k8s.io/apimachinery/pkg/util/sets"
	"log"
	"math"
	"os"

	_ "image/png"
)

func RgbToGray(c color.RGBA) color.Gray {
	gray := 0.3*float64(c.R) + 0.59*float64(c.G) + 0.11*float64(c.B)

	var g color.Gray
	g.Y = uint8(gray)
	return g
}

func GrayscaleImage(img image.Image) *image.Gray {
	bounds := img.Bounds()

	r := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(r, r.Bounds(), img, bounds.Min, draw.Src)

	g := image.NewGray(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c, ok := r.At(x, y).(color.RGBA)
			if !ok {
				log.Fatal()
			}
			g.Set(x, y, RgbToGray(c))
		}
	}

	return g
}

type GaussianKernel struct {
	Size          int
	InvNormFactor int
	Factors       [][]int
}

func CreateGaussianKernel() *GaussianKernel {
	var gk GaussianKernel
	gk.Size = 5
	gk.InvNormFactor = 159
	gk.Factors = [][]int{
		[]int{2, 4, 5, 4, 2},
		[]int{4, 9, 12, 9, 4},
		[]int{5, 12, 15, 12, 5},
		[]int{4, 9, 12, 9, 4},
		[]int{2, 4, 5, 4, 2},
	}

	return &gk
}

func GaussianBlur(img *image.Gray) *image.Gray {
	gk := CreateGaussianKernel()

	bounds := img.Bounds()
	var pxval int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pxval = 0
			for j := -gk.Size / 2; j <= gk.Size/2; j++ {
				for i := -gk.Size / 2; i <= gk.Size/2; i++ {
					m := x + i
					n := y + j

					if m < bounds.Min.X {
						m = bounds.Min.X
					} else if m >= bounds.Max.X {
						m = bounds.Max.X - 1
					}

					if n < bounds.Min.Y {
						n = bounds.Min.Y
					} else if n >= bounds.Max.Y {
						n = bounds.Max.Y - 1
					}

					pxval += int(img.GrayAt(m, n).Y) * gk.Factors[j+(gk.Size/2)][i+(gk.Size/2)]
				}
			}
			pxval /= gk.InvNormFactor
			img.SetGray(x, y, color.Gray{uint8(pxval)})
		}
	}
	return img
}

type SobelKernel struct {
	Size      int
	Direction string
	Factors   [][]int
}

func CreateSobelKernel(dir string) *SobelKernel {
	var sk SobelKernel
	sk.Size = 3

	switch dir {
	case "X":
		sk.Direction = "X"
		sk.Factors = [][]int{
			[]int{-1, 0, 1},
			[]int{-2, 0, 2},
			[]int{-1, 0, 1},
		}
	case "Y":
		sk.Direction = "Y"
		sk.Factors = [][]int{
			[]int{-1, -2, -1},
			[]int{0, 0, 0},
			[]int{1, 2, 1},
		}
	}

	return &sk
}

// type GradientValuesMatrix [][]int

// func CreateGradientValuesMatrix(x, y int) GradientValuesMatrix {
// 	var gvm GradientValuesMatrix

// 	gvm = make([][]int, y)

// 	for j := 0; j < y; j++ {
// 		gvm[j] = make([]int, x)
// 	}

// 	return gvm
// }

type GradientDirection uint8

const (
	zero GradientDirection = iota
	fortyfive
	ninety
	onethreefive
)

func CalcGradientDirection(x, y int) GradientDirection {
	theta := math.Atan(float64(y) / float64(x))

	var gd GradientDirection

	switch {
	case theta <= math.Pi/8:
		gd = zero
	case theta <= 3*math.Pi/8:
		gd = fortyfive
	case theta <= 5*math.Pi/8:
		gd = ninety
	case theta <= 7*math.Pi/8:
		gd = ninety
	default:
		gd = zero
	}

	return gd
}

type ImageGradients struct {
	Value     [][]int
	Direction [][]GradientDirection
}

func CreateImageGradients(x, y int) *ImageGradients {
	var ig ImageGradients

	ig.Value = make([][]int, y)
	ig.Direction = make([][]GradientDirection, y)

	for j := 0; j < y; j++ {
		ig.Value[j] = make([]int, x)
		ig.Direction[j] = make([]GradientDirection, x)
	}

	return &ig
}

func PixelNonmaxSuppression(ig *ImageGradients, x, y int) {
	var above, below image.Point

	maxX := len(ig.Value[0])
	maxY := len(ig.Value)

	switch ig.Direction[y][x] {
	case zero:
		above.X = x + 1
		above.Y = y
		below.X = x - 1
		below.Y = y
	case fortyfive:
		above.X = x + 1
		above.Y = y + 1
		below.X = x - 1
		below.Y = y - 1
	case ninety:
		above.X = x
		above.Y = y + 1
		below.X = x
		below.Y = y - 1
	case onethreefive:
		above.X = x - 1
		above.Y = y + 1
		below.X = x + 1
		below.Y = y - 1
	}

	if above.X >= 0 && above.X < maxX && above.Y >= 0 && above.Y < maxY && ig.Value[above.Y][above.X] > ig.Value[y][x] {
		ig.Value[y][x] = 0
	} else if below.X >= 0 && below.X < maxX && below.Y >= 0 && below.Y < maxY && ig.Value[below.Y][below.X] > ig.Value[y][x] {
		ig.Value[y][x] = 0
	}
}

func (ig *ImageGradients) NonmaxSuppression() *ImageGradients {
	x := len(ig.Value[0])
	y := len(ig.Value)

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			PixelNonmaxSuppression(ig, i, j)
		}
	}

	return ig
}

func (ig *ImageGradients) NeighbourOverThreshold(x, y, thr int) bool {
	maxX := len(ig.Value[0])
	maxY := len(ig.Value)

	if x-1 >= 0 && ig.Value[y][x-1] >= thr {
		return true
	}
	if x-1 >= 0 && y-1 >= 0 && ig.Value[y-1][x-1] >= thr {
		return true
	}
	if y-1 >= 0 && ig.Value[y-1][x] >= thr {
		return true
	}
	if x+1 < maxX && y-1 >= 0 && ig.Value[y-1][x+1] >= thr {
		return true
	}
	if x+1 < maxX && ig.Value[y][x+1] >= thr {
		return true
	}
	if x+1 < maxX && y+1 < maxY && ig.Value[y+1][x+1] >= thr {
		return true
	}
	if y+1 < maxY && ig.Value[y+1][x] >= thr {
		return true
	}
	if x-1 >= 0 && y+1 < maxY && ig.Value[y+1][x-1] >= thr {
		return true
	}
	return false
}

func (ig *ImageGradients) BasicThresholdSuppression() *ImageGradients {
	x := len(ig.Value[0])
	y := len(ig.Value)

	maxVal := 0
	minVal := 2 ^ 8

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			if ig.Value[j][i] > maxVal {
				maxVal = ig.Value[j][i]
			}
			if ig.Value[j][i] < minVal {
				minVal = ig.Value[j][i]
			}
		}
	}

	fmt.Printf("\nMax gradient value is %v, Min value is %v\n", maxVal, minVal)

	upperThreshold := maxVal * 6 / 10
	lowerThreshold := maxVal * 2 / 10

	fmt.Printf("Upper threshold is %v, lower threshold is %v\n", upperThreshold,
		lowerThreshold)

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			if ig.Value[j][i] >= upperThreshold {
				ig.Value[j][i] = 255
			} else if ig.Value[j][i] >= lowerThreshold && ig.NeighbourOverThreshold(i, j, upperThreshold) {
				ig.Value[j][i] = 255
			} else {
				ig.Value[j][i] = 0
			}
		}
	}

	return ig
}

func (ig *ImageGradients) FollowEdge(x, y, upper, lower int, accEdges sets.Set[image.Point]) {
	maxX := len(ig.Value[0])
	maxY := len(ig.Value)

	if x-1 >= 0 && ig.Value[y][x-1] >= lower && ig.Value[y][x-1] < upper && !accEdges.Has(image.Point{x - 1, y}) {
		accEdges.Insert(image.Point{x - 1, y})
		ig.FollowEdge(x-1, y, upper, lower, accEdges)
	}
	if x-1 >= 0 && y-1 >= 0 && ig.Value[y-1][x-1] >= lower && ig.Value[y-1][x-1] < upper && !accEdges.Has(image.Point{x - 1, y - 1}) {
		accEdges.Insert(image.Point{x - 1, y - 1})
		ig.FollowEdge(x-1, y-1, upper, lower, accEdges)
	}
	if y-1 >= 0 && ig.Value[y-1][x] >= lower && ig.Value[y-1][x] < upper && !accEdges.Has(image.Point{x, y - 1}) {
		accEdges.Insert(image.Point{x, y - 1})
		ig.FollowEdge(x, y-1, upper, lower, accEdges)
	}
	if x+1 < maxX && y-1 >= 0 && ig.Value[y-1][x+1] >= lower && ig.Value[y-1][x+1] < upper && !accEdges.Has(image.Point{x + 1, y - 1}) {
		accEdges.Insert(image.Point{x + 1, y - 1})
		ig.FollowEdge(x+1, y-1, upper, lower, accEdges)
	}
	if x+1 < maxX && ig.Value[y][x+1] >= lower && ig.Value[y][x+1] < upper && !accEdges.Has(image.Point{x + 1, y}) {
		accEdges.Insert(image.Point{x + 1, y})
		ig.FollowEdge(x+1, y, upper, lower, accEdges)
	}
	if x+1 < maxX && y+1 < maxY && ig.Value[y+1][x+1] >= lower && ig.Value[y+1][x+1] < upper && !accEdges.Has(image.Point{x + 1, y + 1}) {
		accEdges.Insert(image.Point{x + 1, y + 1})
		ig.FollowEdge(x+1, y+1, upper, lower, accEdges)
	}
	if y+1 < maxY && ig.Value[y+1][x] >= lower && ig.Value[y+1][x] < upper && !accEdges.Has(image.Point{x, y + 1}) {
		accEdges.Insert(image.Point{x, y + 1})
		ig.FollowEdge(x, y+1, upper, lower, accEdges)
	}
	if x-1 >= 0 && y+1 < maxY && ig.Value[y+1][x-1] >= lower && ig.Value[y+1][x-1] < upper && !accEdges.Has(image.Point{x - 1, y + 1}) {
		accEdges.Insert(image.Point{x - 1, y + 1})
		ig.FollowEdge(x-1, y+1, upper, lower, accEdges)
	}
}

func (ig *ImageGradients) LineFollowingThresholdSuppression() *ImageGradients {
	x := len(ig.Value[0])
	y := len(ig.Value)

	maxVal := 0
	minVal := 2 ^ 8

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			if ig.Value[j][i] > maxVal {
				maxVal = ig.Value[j][i]
			}
			if ig.Value[j][i] < minVal {
				minVal = ig.Value[j][i]
			}
		}
	}

	fmt.Printf("\nMax gradient value is %v, Min value is %v\n", maxVal, minVal)

	upperThreshold := 150
	lowerThreshold := 20

	acceptedEdges := make(sets.Set[image.Point])

	fmt.Printf("Upper threshold is %v, lower threshold is %v\n", upperThreshold,
		lowerThreshold)

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			if ig.Value[j][i] >= upperThreshold {
				acceptedEdges.Insert(image.Point{i, j})
				ig.FollowEdge(i, j, upperThreshold, lowerThreshold, acceptedEdges)
			}
		}
	}

	intensityScaleFactor := 255.0 / float64(maxVal)
	fmt.Printf("Intensity scale factor is %v\n", intensityScaleFactor)

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			if !acceptedEdges.Has(image.Point{i, j}) {
				ig.Value[j][i] = 0
			} else {
				ig.Value[j][i] = int(float64(ig.Value[j][i]) * intensityScaleFactor) * 3 / 4 + 64
			}
		}
	}

	return ig
}

func (ig *ImageGradients) GrayscaleImage() *image.Gray {
	x := len(ig.Value[0])
	y := len(ig.Value)

	gray := image.NewGray(image.Rect(0, 0, x, y))

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			gray.SetGray(i, j, color.Gray{uint8(ig.Value[j][i])})
		}
	}

	return gray
}

func InvertGrayscaleImage(img *image.Gray) *image.Gray {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.SetGray(x, y, color.Gray{255 - img.GrayAt(x, y).Y})
		}
	}
	return img
}

func SobelFilter(img *image.Gray) *ImageGradients {
	imgSize := img.Bounds().Size()

	// gx := CreateGradientValuesMatrix(imgSize.X, imgSize.Y)
	// gy := CreateGradientValuesMatrix(imgSize.X, imgSize.Y)

	skx := CreateSobelKernel("X")
	sky := CreateSobelKernel("Y")

	ig := CreateImageGradients(imgSize.X, imgSize.Y)

	var gxval, gyval, imgval int

	var gradientHistogram [17]int

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			gxval, gyval = 0, 0
			for j := -skx.Size / 2; j <= skx.Size/2; j++ {
				for i := -skx.Size / 2; i <= skx.Size/2; i++ {
					m := x + i
					n := y + j

					if m < img.Bounds().Min.X {
						m = img.Bounds().Min.X
					} else if m >= img.Bounds().Max.X {
						m = img.Bounds().Max.X - 1
					}

					if n < img.Bounds().Min.Y {
						n = img.Bounds().Min.Y
					} else if n >= img.Bounds().Max.Y {
						n = img.Bounds().Max.Y - 1
					}

					imgval = int(img.GrayAt(m, n).Y)
					gxval += imgval * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]
					gyval += imgval * sky.Factors[j+(skx.Size/2)][i+(skx.Size/2)]
				}
			}
			// gx[y][x] = gxval
			// gy[y][x] = gyval
			// fmt.Printf("gx, gy = %v, %v", gx[y][x], *gy[y][x])

			ig.Value[y][x] = int(math.Sqrt(float64(gxval*gxval) + float64(gyval*gyval)))
			ig.Direction[y][x] = CalcGradientDirection(gxval, gyval)

			if ig.Value[y][x] == 0 {
				gradientHistogram[16]++
			} else {
				gradientHistogram[ig.Value[y][x]>>6]++
			}
		}
	}

	fmt.Print("\nGradient values:\n")
	fmt.Printf("%-12s   %-8s\n", "Gradient val", "Count")
	fmt.Printf("        %04d   %8d\n", 0, gradientHistogram[16])
	for i, h := range gradientHistogram[0:16] {
		fmt.Printf("   %04d-%04d   %8d\n", i<<6, (i+1)<<6-1, h)
	}

	return ig
}

func main() {
	fmt.Println("Hello, I'm cic!")

	fmt.Print("Reading in file...")

	reader, err := os.Open("micawber-bathtime.jpg")
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

	fmt.Print("Converting to grayscale image...")
	grayImg := GrayscaleImage(img)
	fmt.Print(" Done\n")
	fmt.Print("Applying Gaussian blur...")
	grayImg = GaussianBlur(grayImg)
	fmt.Print(" Done\n")
	fmt.Print("Applying Sobel filter...")
	ig := SobelFilter(grayImg)
	fmt.Print(" Done\n")
	fmt.Print("Applying non-max suppression...")
	ig = ig.NonmaxSuppression()
	fmt.Print(" Done\n")
	fmt.Print("Applying threshold suppression...")
	// ig = ig.BasicThresholdSuppression()
	ig = ig.LineFollowingThresholdSuppression()
	fmt.Print(" Done\n")
	fmt.Print("Converting edge gradients to grayscale image...")
	grayImg = ig.GrayscaleImage()
	fmt.Print(" Done\n")
	fmt.Print("Inverting image to make edges black...")
	grayImg = InvertGrayscaleImage(grayImg)
	fmt.Print(" Done\n")

	fmt.Print("Saving to output file, \"edited.jpg\"...")
	outputFile, err := os.Create("edited.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	var imgOptions jpeg.Options
	imgOptions.Quality = 100

	jpeg.Encode(outputFile, grayImg, &imgOptions)
	fmt.Print(" Done\n")
}
