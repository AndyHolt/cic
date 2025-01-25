package cic

import (
	"image/color"
	"testing"
)

var colors = map[string]struct {
	rgba color.RGBA
	hsva HSVA
}{
	"black": {
		rgba: color.RGBA{0, 0, 0, 255},
		hsva: HSVA{0, 0, 0, 255},
	},
	"white": {
		rgba: color.RGBA{255, 255, 255, 255},
		hsva: HSVA{0, 0, 255, 255},
	},
	"red": {
		rgba: color.RGBA{255, 0, 0, 255},
		hsva: HSVA{0, 255, 255, 255},
	},
	"lime": {
		rgba: color.RGBA{0, 255, 0, 255},
		hsva: HSVA{120, 255, 255, 255},
	},
	"blue": {
		rgba: color.RGBA{0, 0, 255, 255},
		hsva: HSVA{240, 255, 255, 255},
	},
	"yellow": {
		rgba: color.RGBA{255, 255, 0, 255},
		hsva: HSVA{60, 255, 255, 255},
	},
	"cyan": {
		rgba: color.RGBA{0, 255, 255, 255},
		hsva: HSVA{180, 255, 255, 255},
	},
	"magenta": {
		rgba: color.RGBA{255, 0, 255, 255},
		hsva: HSVA{300, 255, 255, 255},
	},
	"silver": {
		rgba: color.RGBA{191, 191, 191, 255},
		hsva: HSVA{0, 0, 191, 255},
	},
	"Gray": {
		rgba: color.RGBA{128, 128, 128, 255},
		hsva: HSVA{0, 0, 128, 255},
	},
	"maroon": {
		rgba: color.RGBA{128, 0, 0, 255},
		hsva: HSVA{0, 255, 128, 255},
	},
	"olive": {
		rgba: color.RGBA{128, 128, 0, 255},
		hsva: HSVA{60, 255, 128, 255},
	},
	"green": {
		rgba: color.RGBA{0, 128, 0, 255},
		hsva: HSVA{120, 255, 128, 255},
	},
	"purple": {
		rgba: color.RGBA{128, 0, 128, 255},
		hsva: HSVA{300, 255, 128, 255},
	},
	"teal": {
		rgba: color.RGBA{0, 128, 128, 255},
		hsva: HSVA{180, 255, 128, 255},
	},
	"navy": {
		rgba: color.RGBA{0, 0, 128, 255},
		hsva: HSVA{240, 255, 128, 255},
	},
	"alpha_normalisation/white": {
		rgba: color.RGBA{60, 60, 60, 60},
		hsva: HSVA{0, 0, 255, 60},
	},
	"alpha_normalisation/green": {
		rgba: color.RGBA{0, 108, 0, 215},
		hsva: HSVA{120, 255, 128, 215},
	},
}

func TestHSVA2RGBA(t *testing.T) {
	t.Parallel()
	for name, c := range colors {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, g, b, a := c.hsva.RGBA()
			rgba := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			if rgba != c.rgba {
				t.Fatalf("Conversion of colour %s from HSVA to RGBA failed. Expected %v, got %v\n", name, c.rgba, rgba)
			}
		})
	}
}

func TestRGBA2HSVA(t *testing.T) {
	// t.Parallel()
	for name, c := range colors {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			h, s, v, a := ConvertRGBA2HSVA(c.rgba)
			hsva := HSVA{uint16(h), uint8(s), uint8(v), uint8(a)}
			if c.hsva != hsva {
				t.Fatalf("Conversion of colour %v from RGBA to HSVA failed. Expected %v, got %v\n", name, c.hsva, hsva)
			}
		})
	}
}
