package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/image/bmp"
)

func openBmpImage(imagePath string) (image.Image, error) {
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

func saveBmpImage(img *image.RGBA, filename string) error {
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

func clampUint8(value int) uint8 {
	if value > 255 {
		return 255
	}
	if value < 0 {
		return 0
	}
	return uint8(value)
}

func adjustBrightness(img image.Image, brightness int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	factor := (brightness * 255) / 100

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := int(r >> 8)
			g8 := int(g >> 8)
			b8 := int(b >> 8)

			newR := clampUint8(r8 + factor)
			newG := clampUint8(g8 + factor)
			newB := clampUint8(b8 + factor)

			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <command> <value> <bmp_image_path>")
		return
	}

	command := os.Args[1]
	value := os.Args[2]
	imagePath := os.Args[3]

	img, err := openBmpImage(imagePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	originalName := filepath.Base(imagePath)
	originalNameWithoutExt := originalName[:len(originalName)-len(filepath.Ext(originalName))]

	var outputFileName string
	var newImg *image.RGBA

	switch command {

	case "--brightness":
		brightness, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Brightness value must be int number: %v", err)
		}

		newImg = adjustBrightness(img, brightness)
		outputFileName = fmt.Sprintf("%s_altered_brightness.bmp", originalNameWithoutExt)

	case "--contrast":
		contrast, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Contrast value must be int number: %v", err)
		}

		if contrast < -255 || contrast > 255 {
			log.Fatalf("Contrast value must be in the range of -255 to 255")
		}

		newImg = adjustContrast(img, contrast)
		outputFileName = fmt.Sprintf("%s_altered_contrast.bmp", originalNameWithoutExt)

	default:
		fmt.Println("Unknown commend")
		return
	}

	err = saveBmpImage(newImg, outputFileName)
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	} else {
		fmt.Printf("Image saved successfully as: %s\n", outputFileName)
	}
}
