package manipulations

import (
	"math"
	"testing"
)

// Reference for colors and their HSV values: https://www.rapidtables.com/convert/color/hsv-to-rgb.html

func TestRGBToHSV(t *testing.T) {
	tests := []struct {
		r, g, b   uint8
		expectedH float64
		expectedS float64
		expectedV float64
	}{
		{0, 0, 0, 0, 0, 0},         // Black
		{255, 255, 255, 0, 0, 1},   // White
		{255, 0, 0, 0, 1, 1},       // Red
		{0, 255, 0, 120, 1, 1},     // Green
		{0, 0, 255, 240, 1, 1},     // Blue
		{255, 255, 0, 60, 1, 1},    // Yellow
		{128, 128, 128, 0, 0, 0.5}, // Gray
		{0, 255, 255, 180, 1, 1},   // Cyan
		{255, 0, 255, 300, 1, 1},   // Magenta
		{0, 0, 128, 240, 1, 0.5},   // Teal
	}

	for _, tt := range tests {
		h, s, v := RGBToHSV(tt.r, tt.g, tt.b)
		if !almostEqual(h, tt.expectedH) || !almostEqual(s, tt.expectedS) || !almostEqual(v, tt.expectedV) {
			t.Errorf("RGB(%d, %d, %d) => HSV(%.2f, %.2f, %.2f), expected HSV(%.2f, %.2f, %.2f)",
				tt.r, tt.g, tt.b, h, s, v, tt.expectedH, tt.expectedS, tt.expectedV)
		}
	}
}

func TestHSVToRGB(t *testing.T) {
	tests := []struct {
		h, s, v   float64
		expectedR uint8
		expectedG uint8
		expectedB uint8
	}{
		{0, 0, 0, 0, 0, 0},         // Black
		{0, 0, 1, 255, 255, 255},   // White
		{0, 1, 1, 255, 0, 0},       // Red
		{120, 1, 0.5, 0, 128, 0},   // Green
		{240, 1, 1, 0, 0, 255},     // Blue
		{60, 1, 1, 255, 255, 0},    // Yellow
		{0, 0, 0.5, 128, 128, 128}, // Gray
		{180, 1, 1, 0, 255, 255},   // Cyan
		{300, 1, 1, 255, 0, 255},   // Magenta
		{180, 1, 0.5, 0, 128, 128}, // Teal
	}

	for _, tt := range tests {
		r, g, b := HSVToRGB(tt.h, tt.s, tt.v)
		if uint8(r+0.5) != tt.expectedR || uint8(g+0.5) != tt.expectedG || uint8(b+0.5) != tt.expectedB {
			t.Errorf("HSV(%.2f, %.2f, %.2f) => RGB(%d, %d, %d), expected RGB(%d, %d, %d)",
				tt.h, tt.s, tt.v, uint8(r+0.5), uint8(g+0.5), uint8(b+0.5), tt.expectedR, tt.expectedG, tt.expectedB)
		}
	}
}

func TestRGBToHSVToRGBRoundTrip(t *testing.T) {
	tests := []struct {
		r, g, b uint8
	}{
		{0, 0, 0},       // Black
		{255, 255, 255}, // White
		{255, 0, 0},     // Red
		{0, 255, 0},     // Green
		{0, 0, 255},     // Blue
		{255, 255, 0},   // Yellow
		{127, 127, 127}, // Gray
		{0, 255, 255},   // Cyan
		{255, 0, 255},   // Magenta
		{0, 0, 128},     // Teal
	}

	for _, tt := range tests {
		h, s, v := RGBToHSV(tt.r, tt.g, tt.b)
		r, g, b := HSVToRGB(h, s, v)
		if uint8(r+0.5) != tt.r || uint8(g+0.5) != tt.g || uint8(b+0.5) != tt.b {
			t.Errorf("RGB(%d, %d, %d) => HSV(%.2f, %.2f, %.2f) => RGB(%d, %d, %d), expected RGB(%d, %d, %d)",
				tt.r, tt.g, tt.b, h, s, v, uint8(r+0.5), uint8(g+0.5), uint8(b+0.5), tt.r, tt.g, tt.b)
		}
	}
}

// almostEqual compares two float64 values and returns true if they are within a small epsilon range.
func almostEqual(a, b float64) bool {
	const epsilon = 0.01
	return math.Abs(a-b) < epsilon
}
