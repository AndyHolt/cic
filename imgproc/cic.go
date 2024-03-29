// Next steps:
// 1. DONE Tidy up code.
// 2. DONE Store size in ImageGradients datastructure so don't need to compute
//    repeatedly
// 3. DONE Make parameters selectable (thresholds, Guassian blur size & sigma,
//    whether or not to use Gaussian blur),
// 4. DONE Make CLI interface
// 5. Make a GUI/web interface to allow fast experimentation with parameter
//    tuning

package cic

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
	X         int
	Y         int
}

func CreateImageGradients(x, y int) *ImageGradients {
	var ig ImageGradients

	ig.X, ig.Y = x, y

	ig.Value = make([][]int, y)
	ig.Direction = make([][]GradientDirection, y)

	for j := 0; j < y; j++ {
		ig.Value[j] = make([]int, x)
		ig.Direction[j] = make([]GradientDirection, x)
	}

	return &ig
}

// Returns True if a pixel should be retained during non-max suppression, and
// False if it should be discarded.
func PixelNonmaxSuppression(ig *ImageGradients, x, y, distance int) bool {
	var above, below image.Point

	for d := 1; d <= distance; d++ {
		switch ig.Direction[y][x] {
		case zero:
			above.X = x + d
			above.Y = y
			below.X = x - d
			below.Y = y
		case fortyfive:
			above.X = x + d
			above.Y = y + d
			below.X = x - d
			below.Y = y - d
		case ninety:
			above.X = x
			above.Y = y + d
			below.X = x
			below.Y = y - d
		case onethreefive:
			above.X = x - d
			above.Y = y + d
			below.X = x + d
			below.Y = y - d
		}

		if above.X >= 0 && above.X < ig.X && above.Y >= 0 && above.Y < ig.Y && ig.Value[above.Y][above.X] > ig.Value[y][x] {
			return false
		} else if below.X >= 0 && below.X < ig.X && below.Y >= 0 && below.Y < ig.Y && ig.Value[below.Y][below.X] > ig.Value[y][x] {
			return false
		}
	}
	return true
}

func (ig *ImageGradients) NonmaxSuppression(distance int) *ImageGradients {
	var pixelState [][]bool
	pixelState = make([][]bool, ig.Y)

	for j := 0; j < ig.Y; j++ {
		pixelState[j] = make([]bool, ig.X)
		for i := 0; i < ig.X; i++ {
			pixelState[j][i] = PixelNonmaxSuppression(ig, i, j, distance)
		}
	}

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
			if !pixelState[j][i] {
				ig.Value[j][i] = 0
			}
		}
	}

	return ig
}

func (ig *ImageGradients) NeighbourOverThreshold(x, y, thr int) bool {
	if x-1 >= 0 && ig.Value[y][x-1] >= thr {
		return true
	}
	if x-1 >= 0 && y-1 >= 0 && ig.Value[y-1][x-1] >= thr {
		return true
	}
	if y-1 >= 0 && ig.Value[y-1][x] >= thr {
		return true
	}
	if x+1 < ig.X && y-1 >= 0 && ig.Value[y-1][x+1] >= thr {
		return true
	}
	if x+1 < ig.X && ig.Value[y][x+1] >= thr {
		return true
	}
	if x+1 < ig.X && y+1 < ig.Y && ig.Value[y+1][x+1] >= thr {
		return true
	}
	if y+1 < ig.Y && ig.Value[y+1][x] >= thr {
		return true
	}
	if x-1 >= 0 && y+1 < ig.Y && ig.Value[y+1][x-1] >= thr {
		return true
	}
	return false
}

func (ig *ImageGradients) BasicThresholdSuppression() *ImageGradients {
	maxVal := 0
	minVal := 2 ^ 8

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
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

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
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
	if x+1 < ig.X && y-1 >= 0 && ig.Value[y-1][x+1] >= lower && ig.Value[y-1][x+1] < upper && !accEdges.Has(image.Point{x + 1, y - 1}) {
		accEdges.Insert(image.Point{x + 1, y - 1})
		ig.FollowEdge(x+1, y-1, upper, lower, accEdges)
	}
	if x+1 < ig.X && ig.Value[y][x+1] >= lower && ig.Value[y][x+1] < upper && !accEdges.Has(image.Point{x + 1, y}) {
		accEdges.Insert(image.Point{x + 1, y})
		ig.FollowEdge(x+1, y, upper, lower, accEdges)
	}
	if x+1 < ig.X && y+1 < ig.Y && ig.Value[y+1][x+1] >= lower && ig.Value[y+1][x+1] < upper && !accEdges.Has(image.Point{x + 1, y + 1}) {
		accEdges.Insert(image.Point{x + 1, y + 1})
		ig.FollowEdge(x+1, y+1, upper, lower, accEdges)
	}
	if y+1 < ig.Y && ig.Value[y+1][x] >= lower && ig.Value[y+1][x] < upper && !accEdges.Has(image.Point{x, y + 1}) {
		accEdges.Insert(image.Point{x, y + 1})
		ig.FollowEdge(x, y+1, upper, lower, accEdges)
	}
	if x-1 >= 0 && y+1 < ig.Y && ig.Value[y+1][x-1] >= lower && ig.Value[y+1][x-1] < upper && !accEdges.Has(image.Point{x - 1, y + 1}) {
		accEdges.Insert(image.Point{x - 1, y + 1})
		ig.FollowEdge(x-1, y+1, upper, lower, accEdges)
	}
}

func (ig *ImageGradients) LineFollowingThresholdSuppression(upperThreshold, lowerThreshold int) *ImageGradients {
	maxVal := 0
	minVal := 2 ^ 8

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
			if ig.Value[j][i] > maxVal {
				maxVal = ig.Value[j][i]
			}
			if ig.Value[j][i] < minVal {
				minVal = ig.Value[j][i]
			}
		}
	}

	fmt.Printf("\nMax gradient value is %v, Min value is %v\n", maxVal, minVal)

	acceptedEdges := make(sets.Set[image.Point])

	fmt.Printf("Upper threshold is %v, lower threshold is %v\n", upperThreshold, lowerThreshold)

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
			if ig.Value[j][i] >= upperThreshold {
				acceptedEdges.Insert(image.Point{i, j})
				ig.FollowEdge(i, j, upperThreshold, lowerThreshold, acceptedEdges)
			}
		}
	}

	intensityScaleFactor := 255.0 / float64(maxVal)
	fmt.Printf("Intensity scale factor is %v\n", intensityScaleFactor)

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
			if !acceptedEdges.Has(image.Point{i, j}) {
				ig.Value[j][i] = 0
			} else {
				ig.Value[j][i] = int(float64(ig.Value[j][i])*intensityScaleFactor)*3/4 + 64
			}
		}
	}

	return ig
}

func (ig *ImageGradients) GrayscaleImage() *image.Gray {
	gray := image.NewGray(image.Rect(0, 0, ig.X, ig.Y))

	for j := 0; j < ig.Y; j++ {
		for i := 0; i < ig.X; i++ {
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

func ThickenLineAtPixel(img *image.Gray, x, y, thickness int, col color.Gray) {
	for j := -thickness / 2; j <= thickness/2; j++ {
		for i := -thickness / 2; i <= thickness/2; i++ {
			if img.GrayAt(x+i, y+j).Y > col.Y {
				img.SetGray(x+i, y+j, col)
			}
		}
	}
}

func ThickenLinesByDarkness(
	src *image.Gray,
	thickerThreshold uint8,
	thinnerThreshold uint8,
) *image.Gray {
	bounds := src.Bounds()

	thickEdge := 5
	thinEdge := 3

	dst := image.NewGray(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if src.GrayAt(x, y).Y <= thickerThreshold {
				ThickenLineAtPixel(dst, x, y, thickEdge, src.GrayAt(x, y))
			} else if src.GrayAt(x, y).Y <= thinnerThreshold {
				ThickenLineAtPixel(dst, x, y, thinEdge, src.GrayAt(x, y))
			}
		}
	}

	return dst
}

func SobelFilter(img *image.Gray) *ImageGradients {
	imgSize := img.Bounds().Size()

	skx := CreateSobelKernel("X")
	sky := CreateSobelKernel("Y")

	ig := CreateImageGradients(imgSize.X, imgSize.Y)

	var gxval, gyval, imgval int

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

			ig.Value[y][x] = int(math.Sqrt(float64(gxval*gxval) + float64(gyval*gyval)))
			ig.Direction[y][x] = CalcGradientDirection(gxval, gyval)

		}
	}

	return ig
}

func ConvertImageToColouring(
	filename string,
	outputFilename string,
	sigma float64,
	upperThreshold int,
	lowerThreshold int,
	nonMaxSuppDist int,
	thickerThreshold int,
	thinnerThreshold int,
) {
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

	fmt.Print("Converting to grayscale image...")
	grayImg := GrayscaleImage(img)
	fmt.Print(" Done\n")
	fmt.Print("Applying Gaussian blur...")
	grayImg = GaussianBlur(grayImg, sigma)
	fmt.Print(" Done\n")
	fmt.Print("Applying Sobel filter...")
	ig := SobelFilter(grayImg)
	fmt.Print(" Done\n")
	fmt.Print("Applying non-max suppression...")
	ig = ig.NonmaxSuppression(nonMaxSuppDist)
	fmt.Print(" Done\n")
	fmt.Print("Applying threshold suppression...")
	// ig = ig.BasicThresholdSuppression()
	ig = ig.LineFollowingThresholdSuppression(upperThreshold, lowerThreshold)
	fmt.Print(" Done\n")
	fmt.Print("Converting edge gradients to grayscale image...")
	grayImg = ig.GrayscaleImage()
	fmt.Print(" Done\n")
	fmt.Print("Inverting image to make edges black...")
	grayImg = InvertGrayscaleImage(grayImg)
	fmt.Print(" Done\n")

	fmt.Print("Thickening lines based on threshold values:\n")
	fmt.Printf(
		"Thicker threshold level: %v\nThinner threshold level: %v\n",
		thickerThreshold,
		thinnerThreshold,
	)
	grayImg = ThickenLinesByDarkness(
		grayImg,
		uint8(thickerThreshold),
		uint8(thinnerThreshold),
	)
	fmt.Print("Done\n")

	fmt.Printf("Saving to output file, \"%v\"...", outputFilename)
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	var imgOptions jpeg.Options
	imgOptions.Quality = 100

	jpeg.Encode(outputFile, grayImg, &imgOptions)
	fmt.Print(" Done\n")
}

func main() {
	fmt.Println("Hello, I'm cic!")
	ConvertImageToColouring(
		"./images/postmanpathelicopter.jpg",
		"./images/postmanpathelicopter-colouring.jpg",
		1.0,
		100,
		20,
		1,
		50,
		150,
	)
}
