package cic

import (
	"image/color"
	"math"
)

type HSVA struct {
	H       uint16
	S, V, A uint8
}

func (c HSVA) RGBA() (r, g, b, a uint32) {
	// S & V have values up to 255, encoding values between 0-1
	var sNorm float64 = float64(c.S) / 255
	var vNorm float64 = float64(c.V) / 255

	var cVal float64 = sNorm * vNorm

	var xFactor float64 = math.Mod(float64(c.H)/60.0, 2) - 1
	if xFactor < 0 {
		xFactor *= -1
	}
	var x float64 = cVal * (1 - xFactor)

	var m float64 = (float64(c.V) / 255) - cVal

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
	r, g, b = uint32(rNorm*aFloat), uint32(gNorm*aFloat), uint32(bNorm*aFloat)

	return r, g, b, uint32(c.A)
}

func ConvertRGBA2HSVA(c color.RGBA) (h, s, v, a uint32) {
	// Scale R, G & B values to alpha un-premultiplied values. Then convert R, G & B values to range 0-1.
	floatAlpha := float64(c.A)
	var rDash float64 = float64(c.R) / floatAlpha
	var gDash float64 = float64(c.G) / floatAlpha
	var bDash float64 = float64(c.B) / floatAlpha

	cMax := math.Max(rDash, math.Max(gDash, bDash))
	cMin := math.Min(rDash, math.Min(gDash, bDash))
	delta := cMax - cMin

	var hDash, sDash, vDash float64
	switch {
	case delta == 0.0:
		hDash = 0
	case cMax == rDash:
		hDash = math.Mod((gDash-bDash)/delta, 6) * 60
		// Normalise -ve values to range 0-360 degrees
		if hDash < 0 {
			hDash = 360 + hDash
		}
	case cMax == gDash:
		hDash = (((bDash - rDash) / delta) + 2) * 60
	case cMax == bDash:
		hDash = (((rDash - gDash) / delta) + 4) * 60
	}

	if cMax == 0 {
		sDash = 0
	} else {
		sDash = delta / cMax
	}

	vDash = cMax

	h = uint32(math.Round(hDash))
	s = uint32(sDash * 255)
	v = uint32(vDash * 255)
	a = uint32(c.A)

	return h, s, v, a
}
