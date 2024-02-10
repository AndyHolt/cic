package cic

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"

	_ "image/png"
)

func imageToRGBA(img image.Image) *image.RGBA {
	// check if image already satisfies RGBA type
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba

}

func GaussianBlurColour(img *image.RGBA, sigma float64) *image.RGBA {
	dgk := CreateDiscreteGaussianKernel(sigma)

	bounds := img.Bounds()
	horizBlurImg := image.NewRGBA(bounds)
	pxval := make([]float64, 4, 4)

	// first pass: along rows
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for idx := range pxval {
				pxval[idx] = 0.0
			}
			for i := -dgk.Size / 2; i <= dgk.Size/2; i++ {
				m := x + i

				if m < bounds.Min.X {
					m = bounds.Min.X
				} else if m >= bounds.Max.X {
					m = bounds.Max.X - 1
				}

				pxval[0] += float64(img.RGBAAt(m, y).R) * dgk.Elements[i+(dgk.Size/2)]
				pxval[1] += float64(img.RGBAAt(m, y).G) * dgk.Elements[i+(dgk.Size/2)]
				pxval[2] += float64(img.RGBAAt(m, y).B) * dgk.Elements[i+(dgk.Size/2)]
				pxval[3] += float64(img.RGBAAt(m, y).A) * dgk.Elements[i+(dgk.Size/2)]
			}
			horizBlurImg.SetRGBA(x, y, color.RGBA{
				uint8(pxval[0] / dgk.ScalingFactor),
				uint8(pxval[1] / dgk.ScalingFactor),
				uint8(pxval[2] / dgk.ScalingFactor),
				uint8(pxval[3] / dgk.ScalingFactor),
			})
		}
	}

	// second pass, down columns
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for idx := range pxval {
				pxval[idx] = 0.0
			}
			for j := -dgk.Size / 2; j <= dgk.Size/2; j++ {
				n := y + j

				if n < bounds.Min.Y {
					n = bounds.Min.Y
				} else if n >= bounds.Max.Y {
					n = bounds.Max.Y - 1
				}

				pxval[0] += float64(horizBlurImg.RGBAAt(x, n).R) * dgk.Elements[j+(dgk.Size/2)]
				pxval[1] += float64(horizBlurImg.RGBAAt(x, n).G) * dgk.Elements[j+(dgk.Size/2)]
				pxval[2] += float64(horizBlurImg.RGBAAt(x, n).B) * dgk.Elements[j+(dgk.Size/2)]
				pxval[3] += float64(horizBlurImg.RGBAAt(x, n).A) * dgk.Elements[j+(dgk.Size/2)]
			}
			img.SetRGBA(x, y, color.RGBA{
				uint8(pxval[0] / dgk.ScalingFactor),
				uint8(pxval[1] / dgk.ScalingFactor),
				uint8(pxval[2] / dgk.ScalingFactor),
				uint8(pxval[3] / dgk.ScalingFactor),
			})
		}
	}

	return img
}

func ColourSobelFilter(img *image.RGBA) *ImageGradients {
	imgSize := img.Bounds().Size()

	skx := CreateSobelKernel("X")
	sky := CreateSobelKernel("Y")

	fmt.Printf("X direction Sobel Kernel is: %v\n", skx.Factors)
	fmt.Printf("Y direction Sobel Kernel is: %v\n", sky.Factors)

	ig := CreateImageGradients(imgSize.X, imgSize.Y)

	gx := make([]int, 3, 3)
	gy := make([]int, 3, 3)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// Reset values of gx and gy for each pixel
			for idx := 0; idx < 3; idx++ {
				gx[idx] = 0.0
				gy[idx] = 0.0
			}
			for j := -sky.Size / 2; j <= sky.Size/2; j++ {
				n := y + j

				if n < img.Bounds().Min.Y {
					n = img.Bounds().Min.Y
				} else if n >= img.Bounds().Max.Y {
					n = img.Bounds().Max.Y - 1
				}

				for i := -skx.Size / 2; i <= skx.Size/2; i++ {
					m := x + i

					if m < img.Bounds().Min.X {
						m = img.Bounds().Min.X
					} else if m >= img.Bounds().Max.X {
						m = img.Bounds().Max.X - 1
					}

					// get pixel value at (m, n) for easy reference
					px := img.RGBAAt(m, n)

					// Horizontal edge components
					gx[0] += int(px.R) * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]
					gx[1] += int(px.G) * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]
					gx[2] += int(px.B) * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]

					// Vertical edge components
					gy[0] += int(px.R) * sky.Factors[j+(sky.Size/2)][i+(sky.Size/2)]
					gy[1] += int(px.G) * sky.Factors[j+(sky.Size/2)][i+(sky.Size/2)]
					gy[2] += int(px.B) * sky.Factors[j+(sky.Size/2)][i+(sky.Size/2)]
				}
			}

			// Calculate single gx and gy values for gradient, based on RGB
			// channels together
			// [todo] - is the division by 9 right here? Is that the correct
			// normalisation factor?
			gradX := math.Sqrt(float64((gx[0] * gx[0]) + (gx[1] * gx[1]) + (gx[2] * gx[2])))
			gradY := math.Sqrt(float64((gy[0] * gy[0]) + (gy[1] * gy[1]) + (gy[2] * gy[2])))

			// gradX := (math.Abs(float64(gx[0])) + math.Abs(float64(gx[1])) +
			// math.Abs(float64(gx[2]))) / 3
			// gradY := (math.Abs(float64(gy[0])) + math.Abs(float64(gy[1])) +
			// math.Abs(float64(gy[2]))) / 3

			// gradX := float64(max(gx[0], gx[1], gx[2]))
			// gradY := float64(max(gy[0], gy[1], gy[2]))

			ig.Value[y][x] = int(math.Sqrt((gradX * gradX) + (gradY * gradY)))
			ig.Direction[y][x] = CalcGradientDirection(int(gradX), int(gradY))
		}
	}

	return ig
}

func SeparateColourSobelFilter(img *image.RGBA) (
	*ImageGradients, *ImageGradients, *ImageGradients) {
	imgSize := img.Bounds().Size()

	skx := CreateSobelKernel("X")
	sky := CreateSobelKernel("Y")

	fmt.Printf("X direction Sobel Kernel is: %v\n", skx.Factors)
	fmt.Printf("Y direction Sobel Kernel is: %v\n", sky.Factors)

	igR := CreateImageGradients(imgSize.X, imgSize.Y)
	igG := CreateImageGradients(imgSize.X, imgSize.Y)
	igB := CreateImageGradients(imgSize.X, imgSize.Y)

	gx := make([]int, 3, 3)
	gy := make([]int, 3, 3)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			// Reset values of gx and gy for each pixel
			for idx := 0; idx < 3; idx++ {
				gx[idx] = 0.0
				gy[idx] = 0.0
			}
			for j := -sky.Size / 2; j <= sky.Size/2; j++ {
				n := y + j

				if n < img.Bounds().Min.Y {
					n = img.Bounds().Min.Y
				} else if n >= img.Bounds().Max.Y {
					n = img.Bounds().Max.Y - 1
				}

				for i := -skx.Size / 2; i <= skx.Size/2; i++ {
					m := x + i

					if m < img.Bounds().Min.X {
						m = img.Bounds().Min.X
					} else if m >= img.Bounds().Max.X {
						m = img.Bounds().Max.X - 1
					}

					// get pixel value at (m, n) for easy reference
					px := img.RGBAAt(m, n)

					// Horizontal edge components
					gx[0] += int(px.R) * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]
					gx[1] += int(px.G) * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]
					gx[2] += int(px.B) * skx.Factors[j+(skx.Size/2)][i+(skx.Size/2)]

					// Vertical edge components
					gy[0] += int(px.R) * sky.Factors[j+(sky.Size/2)][i+(sky.Size/2)]
					gy[1] += int(px.G) * sky.Factors[j+(sky.Size/2)][i+(sky.Size/2)]
					gy[2] += int(px.B) * sky.Factors[j+(sky.Size/2)][i+(sky.Size/2)]
				}
			}

			// Keep colours separate to display each colour's edge contribution
			igR.Value[y][x] = int(math.Sqrt(float64(gx[0]*gx[0]) + float64(gy[0]*gy[0])))
			igR.Direction[y][x] = CalcGradientDirection(gx[0], gy[0])

			igG.Value[y][x] = int(math.Sqrt(float64(gx[1]*gx[1]) + float64(gy[1]*gy[1])))
			igG.Direction[y][x] = CalcGradientDirection(gx[1], gy[1])

			igB.Value[y][x] = int(math.Sqrt(float64(gx[2]*gx[2]) + float64(gy[2]*gy[2])))
			igB.Direction[y][x] = CalcGradientDirection(gx[2], gy[2])

		}
	}

	return igR, igG, igB
}

func RunSeparateColourImageProc() {
	filename := "./images/engines.jpg"
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

	fmt.Print("Applying Gaussian blur...")
	blurImg := GaussianBlurColour(rgba, 2.0)
	fmt.Print(" Done\n")

	fmt.Print("Applying Sobel Filter for each colour...")
	igR, igG, igB := SeparateColourSobelFilter(blurImg)
	fmt.Print(" Done\n")

	// fmt.Print("Applying non-max suppression...")
	// ig = ig.NonmaxSuppression()
	// fmt.Print(" Done\n")

	// fmt.Print("Applying threshold suppression...")
	// // ig = ig.BasicThresholdSuppression()
	// ig = ig.LineFollowingThresholdSuppression(800, 400)
	// fmt.Print(" Done\n")

	fmt.Print("Converting edge gradients to grayscale image...")
	grayImgR := igR.GrayscaleImage()
	grayImgG := igG.GrayscaleImage()
	grayImgB := igB.GrayscaleImage()
	fmt.Print(" Done\n")

	fmt.Print("Inverting image to make edges black...")
	grayImgR = InvertGrayscaleImage(grayImgR)
	grayImgG = InvertGrayscaleImage(grayImgG)
	grayImgB = InvertGrayscaleImage(grayImgB)
	fmt.Print(" Done\n")

	fmt.Print("Saving to output files: \"./images/engines-[red|green|blue].jpg\"")
	outputFileRed, err := os.Create("./images/engines-red.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFileRed.Close()

	var imgOptions jpeg.Options
	imgOptions.Quality = 100

	jpeg.Encode(outputFileRed, grayImgR, &imgOptions)

	outputFileGreen, err := os.Create("./images/engines-green.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFileGreen.Close()

	jpeg.Encode(outputFileGreen, grayImgG, &imgOptions)

	outputFileBlue, err := os.Create("./images/engines-blue.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFileBlue.Close()

	jpeg.Encode(outputFileBlue, grayImgB, &imgOptions)
	fmt.Print(" Done\n")
}

// [todo] -- add root app CLI options to colour processing
func RunColourImageProc(filename string, outputFilename string) {
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

	fmt.Print("Applying Gaussian blur...")
	blurImg := GaussianBlurColour(rgba, 2.0)
	fmt.Print(" Done\n")

	fmt.Print("Applying Sobel Filter for each colour...")
	ig := ColourSobelFilter(blurImg)
	fmt.Print(" Done\n")

	fmt.Print("Applying non-max suppression...")
	// [todo] - Replace function parameter with distance value
	ig = ig.NonmaxSuppression(1)
	fmt.Print(" Done\n")

	fmt.Print("Applying threshold suppression...")
	// ig = ig.BasicThresholdSuppression()
	ig = ig.LineFollowingThresholdSuppression(75, 25)
	fmt.Print(" Done\n")

	fmt.Print("Converting edge gradients to grayscale image...")
	grayImg := ig.GrayscaleImage()
	fmt.Print(" Done\n")

	fmt.Print("Inverting image to make edges black...")
	grayImg = InvertGrayscaleImage(grayImg)
	fmt.Print(" Done\n")

	fmt.Printf("Saving to output file: \"%v\"", outputFilename)
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
