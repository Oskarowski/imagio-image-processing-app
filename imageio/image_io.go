package imageio

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/image/bmp"
)

func OpenBmpImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err == nil {
		return img, nil
	}
	log.Default().Printf("Error decoding BMP image file: %v", err)

	// Reset file pointer to the beginning for other formats
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("error resetting file pointer: %v", err)
	}

	// Fallback to generic image decoding
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image file (%s): %v", format, err)
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

func LoadMonochromeBMP(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read BMP header
	header := make([]byte, 54) // Standard BMP header size
	_, err = file.Read(header)
	if err != nil {
		return nil, fmt.Errorf("error reading BMP header: %v", err)
	}

	// Parse dimensions
	width := int(binary.LittleEndian.Uint32(header[18:22]))
	height := int(binary.LittleEndian.Uint32(header[22:26]))
	bitDepth := uint16(binary.LittleEndian.Uint16(header[28:30]))

	if bitDepth != 1 {
		return nil, fmt.Errorf("unsupported bit depth: %d", bitDepth)
	}

	// Calculate padding (each row is padded to 4 bytes)
	rowBytes := (width + 7) / 8 // Each pixel is 1 bit
	rowPadding := (4 - (rowBytes % 4)) % 4

	// Read pixel data
	pixelDataOffset := binary.LittleEndian.Uint32(header[10:14])
	_, err = file.Seek(int64(pixelDataOffset), 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking to pixel data: %v", err)
	}

	pixelData := make([]byte, (rowBytes+rowPadding)*height)
	_, err = file.Read(pixelData)
	if err != nil {
		return nil, fmt.Errorf("error reading pixel data: %v", err)
	}

	img := image.NewGray(image.Rect(0, 0, width, height))

	// Convert BMP pixel data to image.Gray
	for y := 0; y < height; y++ {
		rowStart := y * (rowBytes + rowPadding)
		for x := 0; x < width; x++ {
			byteIndex := rowStart + x/8
			bitIndex := 7 - (x % 8)
			if (pixelData[byteIndex]>>bitIndex)&1 == 1 {
				img.SetGray(x, height-y-1, color.Gray{Y: 255}) // BMP stores rows bottom-to-top
			} else {
				img.SetGray(x, height-y-1, color.Gray{Y: 0})
			}
		}
	}

	return img, nil
}
