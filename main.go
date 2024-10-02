package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
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

	factor := float64(brightness) / 100.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8

			newR := clampUint8(int(float64(r8) + (255 * factor)))
			newG := clampUint8(int(float64(g8) + (255 * factor)))
			newB := clampUint8(int(float64(b8) + (255 * factor)))

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

	switch command {
	case "--brightness":
		brightness, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Brightness value must be number: %v", err)
		}

		newImg := adjustBrightness(img, brightness)
		err = saveBmpImage(newImg, "output.bmp")
		if err != nil {
			log.Fatalf("Error saving file: %v", err)
		} else {
			fmt.Printf("Brightness adjusted successfully and saved to: %s", "output.bmp")
		}
	default:
		fmt.Println("Unknown commend")
	}
}
