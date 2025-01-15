package orthogonal_transforms

import (
	"image"
	"image/color"
	"math"
)

func ConvertImageToComplex(img image.Image) [][]complex128 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	complexImg := make([][]complex128, height)

	for y := 0; y < height; y++ {
		complexImg[y] = make([]complex128, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
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

			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
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
			intensity := uint8(math.Min(math.Max(matrix[y][x], 0), 255))

			img.Set(x, y, color.RGBA{intensity, intensity, intensity, 255})
		}
	}

	return img
}

func ConvertComplexToImage(complexImg [][]complex128) *image.RGBA {
	height := len(complexImg)
	width := len(complexImg[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var min, max float64
	min = math.Inf(1)
	max = math.Inf(-1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := real(complexImg[y][x])
			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
	}

	gamma := 0.5
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := real(complexImg[y][x])
			normalized := (value - min) / (max - min) * 255.0
			corrected := math.Pow(normalized/255.0, 1.0/gamma) * 255.0
			pixelValue := uint8(math.Round(corrected))
			img.Set(x, y, color.RGBA{pixelValue, pixelValue, pixelValue, 255})
		}
	}

	return img
}
