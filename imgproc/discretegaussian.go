package cic

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

func intAbs(n int) int {
	if n >= 0 {
		return n
	} else {
		return -n
	}
}

// t of Discrete Gaussian is σ² of continuous Gaussian, i.e. the variance
// (square of standard deviation)
func DiscreteGaussian(n int, t float64) float64 {
	return math.Exp(-t) * ModBesselIn(intAbs(n), t)
}

type DiscreteGaussianKernel1D struct {
	Size          int
	Variance      float64
	ScalingFactor float64
	Elements      []float64
}

func (dgk DiscreteGaussianKernel1D) String() string {
	s := fmt.Sprintf("Discrete Gaussian Kernel with variance t = %.1f (std dev σ = %.1f)\n", dgk.Variance, math.Sqrt(dgk.Variance))
	s += fmt.Sprintf("Kernel size is: %v\n", dgk.Size)
	s += fmt.Sprintf("Kernel elements are: [")
	for _, e := range dgk.Elements[:dgk.Size-1] {
		s += fmt.Sprintf("%.5f, ", e)
	}
	s += fmt.Sprintf("%.5f]\n", dgk.Elements[dgk.Size-1])
	s += fmt.Sprintf("Applying this DGK will give a scaling factor of %.2f\n", dgk.ScalingFactor)

	return s
}

func sigma2size(sigma float64) int {
	var size int
	size = int(math.Ceil((4 * sigma) + 1))
	if size%2 == 0 {
		return size + 1
	} else {
		return size
	}
}

func CreateDiscreteGaussianKernel(sigma float64) *DiscreteGaussianKernel1D {
	var dgk DiscreteGaussianKernel1D
	dgk.Variance = sigma * sigma
	dgk.Size = sigma2size(sigma)
	dgk.Elements = make([]float64, dgk.Size)

	for i := 0; i <= dgk.Size/2; i++ {
		dgk.Elements[i] = DiscreteGaussian((-dgk.Size/2)+i, dgk.Variance)
		dgk.ScalingFactor += dgk.Elements[i]
		if i != dgk.Size/2 {
			dgk.Elements[dgk.Size-i-1] = dgk.Elements[i]
			dgk.ScalingFactor += dgk.Elements[i]
		}
	}

	return &dgk
}

func GaussianBlur(img *image.Gray, sigma float64) *image.Gray {
	dgk := CreateDiscreteGaussianKernel(sigma)

	bounds := img.Bounds()

	// first pass: along rows
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pxval := 0.0
			for i := -dgk.Size / 2; i <= dgk.Size/2; i++ {
				m := x + i

				if m < bounds.Min.X {
					m = bounds.Min.X
				} else if m >= bounds.Max.X {
					m = bounds.Max.X - 1
				}

				pxval += float64(img.GrayAt(m, y).Y) * dgk.Elements[i+(dgk.Size/2)]
			}
			pxval /= dgk.ScalingFactor
			img.SetGray(x, y, color.Gray{uint8(pxval)})
		}
	}

	// second pass: down columns
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			pxval := 0.0
			for j := -dgk.Size / 2; j <= dgk.Size/2; j++ {
				n := y + j

				if n < bounds.Min.Y {
					n = bounds.Min.Y
				} else if n >= bounds.Max.Y {
					n = bounds.Max.Y - 1
				}

				pxval += float64(img.GrayAt(x, n).Y) * dgk.Elements[j+(dgk.Size/2)]
			}
			pxval /= dgk.ScalingFactor
			img.SetGray(x, y, color.Gray{uint8(pxval)})
		}
	}

	return img
}
