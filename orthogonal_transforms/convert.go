package orthogonal_transforms

import (
	"image"
	"image/color"
)

func ConvertImageToComplex(img image.Image) [][]complex128 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	complexImg := make([][]complex128, height)

	for y := 0; y < height; y++ {
		complexImg[y] = make([]complex128, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			gray := 0.299*float64(r/256) + 0.587*float64(g/256) + 0.114*float64(b/256)
			complexImg[y][x] = complex(gray, 0)

		}
	}

	return complexImg
}

func ConvertToFloatMatrix(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	matrix := make([][]float64, height)

	for y := 0; y < height; y++ {
		matrix[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			gray := 0.299*float64(r/256) + 0.587*float64(g/256) + 0.114*float64(b/256)
			matrix[y][x] = gray

		}
	}

	return matrix
}

func ConvertFloatMatrixToImage(matrix [][]float64) *image.RGBA {
	height := len(matrix)
	width := len(matrix[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			intensity := uint8(matrix[y][x])

			img.Set(x, y, color.RGBA{intensity, intensity, intensity, 255})
		}
	}

	return img
}

func ConvertComplexToImage(complexImg [][]complex128) *image.RGBA {
	height := len(complexImg)
	width := len(complexImg[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixelValue := uint8(real(complexImg[y][x]))
			img.Set(x, y, color.RGBA{pixelValue, pixelValue, pixelValue, 255})
		}
	}

	return img
}
