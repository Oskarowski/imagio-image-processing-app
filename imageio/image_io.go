package imageio

import (
	"fmt"
	"image"
	"os"

	"golang.org/x/image/bmp"
)

func OpenBmpImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding BMP image file: %v", err)
	}

	return img, nil
}

func SaveBmpImage(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating BMP image file: %v", err)
	}
	defer file.Close()

	err = bmp.Encode(file, img)
	if err != nil {
		return fmt.Errorf("error encoding BMP image file: %v", err)
	}

	return nil
}
