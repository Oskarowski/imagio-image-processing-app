package manipulations

import (
	"image"
	"image/color"
	"testing"
)

func generateTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := uint8((x + y) % 256)
			img.Set(x, y, color.Gray{Y: gray})
		}
	}
	return img
}

func BenchmarkApplyConvolutionUniversal(b *testing.B) {
	img := generateTestImage(512, 512)
	mask := [][]int{
		{-1, -1, -1},
		{-1, 9, -1},
		{-1, -1, -1},
	}
	for i := 0; i < b.N; i++ {
		ApplyConvolutionUniversal(img, mask)
	}
}

func BenchmarkApplyConvolutionOptimized(b *testing.B) {
	img := generateTestImage(512, 512)
	for i := 0; i < b.N; i++ {
		ApplyConvolutionOptimized(img)
	}
}
