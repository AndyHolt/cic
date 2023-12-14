package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
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

type GradientValuesMatrix [][]int

func CreateGradientValuesMatrix(x, y int) GradientValuesMatrix {
	var gvm GradientValuesMatrix

	gvm = make([][]int, y)

	for j := 0; j < y; j++ {
		gvm[j] = make([]int, x)
	}

	return gvm
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
		}
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
	imgOptions.Quality = 75

	jpeg.Encode(outputFile, grayImg, &imgOptions)
	fmt.Print(" Done\n")
}
