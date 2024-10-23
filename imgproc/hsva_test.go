package cic

import (
	"image/color"
	"testing"
)

func TestHSVA2RGBABlack(t *testing.T) {
	var black = HSVA{0, 0, 0, 255}
	r, g, b, a := black.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 0, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour black from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAWhite(t *testing.T) {
	var white = HSVA{0, 0, 255, 255}
	r, g, b, a := white.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{255, 255, 255, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour white from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBARed(t *testing.T) {
	var red = HSVA{0, 255, 255, 255}
	r, g, b, a := red.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{255, 0, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour red from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBALime(t *testing.T) {
	var lime = HSVA{120, 255, 255, 255}
	r, g, b, a := lime.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 255, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour lime from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBABlue(t *testing.T) {
	var blue = HSVA{240, 255, 255, 255}
	r, g, b, a := blue.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 0, 255, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour blue from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAYellow(t *testing.T) {
	var yellow = HSVA{60, 255, 255, 255}
	r, g, b, a := yellow.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{255, 255, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour yellow from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBACyan(t *testing.T) {
	var cyan = HSVA{180, 255, 255, 255}
	r, g, b, a := cyan.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 255, 255, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour cyan from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAMagenta(t *testing.T) {
	var magenta = HSVA{300, 255, 255, 255}
	r, g, b, a := magenta.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{255, 0, 255, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour magenta from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBASilver(t *testing.T) {
	var silver = HSVA{0, 0, 191, 255}
	r, g, b, a := silver.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{191, 191, 191, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour silver from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAGray(t *testing.T) {
	var gray = HSVA{0, 0, 128, 255}
	r, g, b, a := gray.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{128, 128, 128, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour gray from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAMaroon(t *testing.T) {
	var maroon = HSVA{0, 255, 128, 255}
	r, g, b, a := maroon.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{128, 0, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour maroon from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAOlive(t *testing.T) {
	var olive = HSVA{60, 255, 128, 255}
	r, g, b, a := olive.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{128, 128, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour olive from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAGreen(t *testing.T) {
	var green = HSVA{120, 255, 128, 255}
	r, g, b, a := green.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 128, 0, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour green from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBAPurple(t *testing.T) {
	var purple = HSVA{300, 255, 128, 255}
	r, g, b, a := purple.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{128, 0, 128, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour purple from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBATeal(t *testing.T) {
	var teal = HSVA{180, 255, 128, 255}
	r, g, b, a := teal.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 128, 128, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour teal from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestHSVA2RGBANavy(t *testing.T) {
	var navy = HSVA{240, 255, 128, 255}
	r, g, b, a := navy.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 0, 128, 255}

	if rgba != expected {
		t.Fatalf("Conversion of colour navy from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestAlphaNormalisationWhite(t *testing.T) {
	var white = HSVA{0, 0, 255, 60}
	r, g, b, a := white.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{60, 60, 60, 60}

	if rgba != expected {
		t.Fatalf("Alpha-normalisation of white from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}

func TestAlphaNormalisationGreen(t *testing.T) {
	var green = HSVA{120, 255, 128, 215}
	r, g, b, a := green.RGBA()
	var rgba = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	var expected = color.RGBA{0, 107, 0, 215}

	if rgba != expected {
		t.Fatalf("Alpha-normalisation of green from HSVA to RGBA failed. Expected %v, got %v\n", expected, rgba)
	}
}




func ConvertRGBA2HSVA(t *testing.T) {
	t.Fatalf("Not yet implemented")
}
