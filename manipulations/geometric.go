package manipulations

import (
	"fmt"
	"image"
)

func HorizontalFlip(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(bounds.Max.X-x-1, y, img.At(x, y))
		}
	}

	return newImg
}

func VerticalFlip(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, bounds.Max.Y-y-1, img.At(x, y))
		}
	}

	return newImg
}

func DiagonalFlip(img image.Image) *image.RGBA {
	// easy way to do it but, why should you make it easy when you can make it difficult?
	// return horizontalFlip(verticalFlip(img))

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	newImg := image.NewRGBA(image.Rect(0, 0, height, width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(width-x-1, height-y-1, img.At(x, y))
		}
	}

	return newImg
}

func ShrinkImage(img image.Image, factor int) (*image.RGBA, error) {
	if factor <= 0 {
		return nil, fmt.Errorf("factor must be greater than 0")
	}

	bounds := img.Bounds()
	newWidth := bounds.Dx() / factor
	newHeight := bounds.Dy() / factor

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			newImg.Set(x, y, img.At(x*factor, y*factor))
		}
	}

	return newImg, nil
}

func EnlargeImage(img image.Image, factor int) (*image.RGBA, error) {
	if factor <= 0 {
		return nil, fmt.Errorf("factor must be greater than 0")
	}

	bounds := img.Bounds()
	newWidth := bounds.Dx() * factor
	newHeight := bounds.Dy() * factor

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			newImg.Set(x, y, img.At(x/factor, y/factor))
		}
	}

	return newImg, nil
}
