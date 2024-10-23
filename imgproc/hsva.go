package cic

import (
	"math"
)

type HSVA struct {
	H uint16
	S, V, A uint8
}

func (c HSVA) RGBA() (r, g, b, a uint32) {
	// S & V have values up to 255, encoding values between 0-1
	var sNorm float64 = float64(c.S) / 255
	var vNorm float64 = float64(c.V) / 255

	var cVal float64 = sNorm * vNorm

	var xFactor float64 = math.Mod(float64(c.H) / 60.0, 2) - 1
	if xFactor < 0 {
		xFactor *= -1
	}
	var x float64 = cVal * (1 - xFactor)

	var m float64 = (float64(c.V) / 255)  - cVal

	var rDash, gDash, bDash float64
	switch {
	case c.H < 60 || c.H == 360:
		rDash, gDash, bDash = cVal, x, 0
	case c.H < 120:
		rDash, gDash, bDash = x, cVal, 0
	case c.H < 180:
		rDash, gDash, bDash = 0, cVal, x
	case c.H < 240:
		rDash, gDash, bDash = 0, x, cVal
	case c.H < 300:
		rDash, gDash, bDash = x, 0, cVal
	case c.H < 360:
		rDash, gDash, bDash = cVal, 0, x
	}

	// Norm values are all in range (0, 1), will be multiplied by 255 to get non
	// alpha-premultiplied RGB values.
	var rNorm, gNorm, bNorm float64 = rDash + m, gDash + m, bDash + m

	aFloat := float64(c.A)
	r, g, b = uint32(rNorm * aFloat), uint32(gNorm * aFloat), uint32(bNorm * aFloat)

	return r, g, b, uint32(c.A)
}
