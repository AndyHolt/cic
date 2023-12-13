package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	fmt.Print("Hello, I'm cic!\n")

	reader, err := os.Open("micawber-bathtime.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	fmt.Printf("Image has bounds: %v\n", bounds)

	m := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(m, m.Bounds(), img, bounds.Min, draw.Src)

	for y := bounds.Min.Y; y < (bounds.Min.Y + 100); y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			m.Set(x, y, color.RGBA{255, 0, 0, 0})
		}
	}

	outputFile, err := os.Create("edited.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	var imgOptions jpeg.Options
	imgOptions.Quality = 75

	jpeg.Encode(outputFile, m, &imgOptions)
}
