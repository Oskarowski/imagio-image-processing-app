package imageio

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

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
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating output directory: %v", err)
		}
	}

	fullPath := filepath.Join(outputDir, filename)

	file, err := os.Create(fullPath)
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
