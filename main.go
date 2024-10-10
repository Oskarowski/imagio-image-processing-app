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

// The function returns a new image with the adjusted contrast.
//
// Parameters:
//
//	img      - The input image to adjust.
//	contrast - Desired level of contrast.
//
// Returns:
//
//	*image.RGBA - A new image with the adjusted contrast.
func adjustContrast(img image.Image, contrast int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Contrast correction factor formula:
	// https://www.dfstudios.co.uk/articles/programming/image-programming-algorithms/image-processing-algorithms-part-5-contrast-adjustment/
	// https://ie.nitk.ac.in/blog/2020/01/19/algorithms-for-adjusting-brightness-and-contrast-of-an-image/
	var contrastCorrectionFactor float64 = (259.0 * float64(contrast+255)) / (255.0 * float64(259-contrast))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			newR := clampUint8(int(contrastCorrectionFactor*(r8-128) + 128))
			newG := clampUint8(int(contrastCorrectionFactor*(g8-128) + 128))
			newB := clampUint8(int(contrastCorrectionFactor*(b8-128) + 128))

			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}

func negativeImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			newR := 255 - uint8(r>>8)
			newG := 255 - uint8(g>>8)
			newB := 255 - uint8(b>>8)

			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}

func horizontalFlip(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(bounds.Max.X-x-1, y, img.At(x, y))
		}
	}

	return newImg
}

func verticalFlip(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(x, bounds.Max.Y-y-1, img.At(x, y))
		}
	}

	return newImg
}

func diagonalFlip(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			newImg.Set(y, x, img.At(x, y))
		}
	}

	return newImg
}

func shrinkImage(img image.Image, factor int) (*image.RGBA, error) {
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

func enlargeImage(img image.Image, factor int) (*image.RGBA, error) {
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <command> [<value>] <bmp_image_path>")
		return
	}

	command := os.Args[1]
	imagePath := os.Args[len(os.Args)-1]

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
		if len(os.Args) < 4 {
			log.Fatalf("Brightness command requires a value.")
		}

		brightness, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Brightness value must be int number: %v", err)
		}

		newImg = adjustBrightness(img, brightness)
		outputFileName = fmt.Sprintf("%s_altered_brightness.bmp", originalNameWithoutExt)

	case "--contrast":
		if len(os.Args) < 4 {
			log.Fatalf("Contrast command requires a value.")
		}

		contrast, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Contrast value must be int number: %v", err)
		}

		if contrast < -255 || contrast > 255 {
			log.Fatalf("Contrast value must be in the range of -255 to 255")
		}

		newImg = adjustContrast(img, contrast)
		outputFileName = fmt.Sprintf("%s_altered_contrast.bmp", originalNameWithoutExt)

	case "--negative":
		newImg = negativeImage(img)
		outputFileName = fmt.Sprintf("%s_negative.bmp", originalNameWithoutExt)

	case "--hflip":
		newImg = horizontalFlip(img)
		outputFileName = fmt.Sprintf("%s_horizontal_flip.bmp", originalNameWithoutExt)

	case "--vflip":
		newImg = verticalFlip(img)
		outputFileName = fmt.Sprintf("%s_vertical_flip.bmp", originalNameWithoutExt)

	case "--dflip":
		newImg = diagonalFlip(img)
		outputFileName = fmt.Sprintf("%s_diagonal_flip.bmp", originalNameWithoutExt)

	case "--shrink":
		if len(os.Args) < 4 {
			log.Fatal("Shrink command requires a shrink factor value")
		}

		factor, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Shrink factor value must be int number: %v", err)
		}

		newImg, err = shrinkImage(img, factor)
		if err != nil {
			log.Fatalf("Error shrinking image: %v", err)
		}

		outputFileName = fmt.Sprintf("%s_shrunk_by_%dx.bmp", originalNameWithoutExt, factor)

	case "--enlarge":
		if len(os.Args) < 4 {
			log.Fatal("Enlarge command requires a enlargement factor value")
		}

		factor, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Enlarge factor value must be int number: %v", err)
		}

		newImg, err = enlargeImage(img, factor)
		if err != nil {
			log.Fatalf("Error enlarging image: %v", err)
		}

		outputFileName = fmt.Sprintf("%s_enlarged_by_%dx.bmp", originalNameWithoutExt, factor)

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
