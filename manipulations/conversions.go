package manipulations

import (
	"image"
	"image/color"
	"math"
)

// RGBToHSV converts RGB color values to the HSV (Hue, Saturation, Value) color model.
//
// Parameters:
// - r, g, b: RGB color values ranging from 0 to 255.
//
// Returns:
// - h: Hue, in degrees [0, 360].
// - s: Saturation, as a fraction [0, 1].
// - v: Value (intensitivity), as a fraction [0, 1].
//
// The algorithm follows a common RGB-to-HSV conversion method.
// Reference: https://math.stackexchange.com/questions/556341/rgb-to-hsv-color-conversion-algorithm
func RGBToHSV(r, g, b uint8) (h, s, v float64) {
	rf, gf, bf := float64(r)/255.0, float64(g)/255.0, float64(b)/255.0
	ColorMax := math.Max(rf, math.Max(gf, bf))
	ColorMin := math.Min(rf, math.Min(gf, bf))
	delta := ColorMax - ColorMin

	v = ColorMax

	if delta == 0 {
		return 0, 0, v
	}

	h, s = 0, 0
	s = delta / ColorMax
	switch ColorMax {
	case rf:
		h = (gf - bf) / delta
	case gf:
		h = 2 + (bf-rf)/delta
	case bf:
		h = 4 + (rf-gf)/delta
	}
	h *= 60
	if h < 0 {
		h += 360
	}

	return h, s, v
}

// HSVToRGB converts HSV (Hue, Saturation, Value) color values to RGB (Red, Green, Blue) color values.
// The input parameters are:
//   - h: Hue, a value between 0 and 360 degrees.
//   - s: Saturation, as a fraction [0, 1].
//   - v: Value (intensitivity), as a fraction [0, 1].
//
// The function returns three float64 values representing the RGB color components, each in the range of 0 to 255.
// Reference: https://cs.stackexchange.com/questions/64549/convert-hsv-to-rgb-colors
func HSVToRGB(h, s, v float64) (r, g, b float64) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	}

	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	return r, g, b
}

func RoundTripToHSVtoRGB(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			h, s, v := RGBToHSV(pixel.R, pixel.G, pixel.B)
			r, g, b := HSVToRGB(h, s, v)

			newImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), pixel.A})

		}
	}

	return newImg
}
