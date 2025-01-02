package morphological

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
)

// Helpful presentation: https://www.ee.nthu.edu.tw/clhuang/09420EE368000DIP/chapter09.pdf

type BinaryImage [][]int

type StructuringElement struct {
	Data    [][]int `json:"data"`
	OriginX int     `json:"originX"`
	OriginY int     `json:"originY"`
}

type StructureElementsJSON struct {
	StructureElements map[string]StructuringElement `json:"structure_elements"`
}

func GetStructureElement(structureElements map[string]StructuringElement, seName string) (StructuringElement, error) {
	if structureElement, exists := structureElements[seName]; exists {
		return structureElement, nil
	}
	return StructuringElement{}, fmt.Errorf("structural Element %s not found", seName)
}

func ConvertIntoBinaryImage(img image.Image) BinaryImage {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	binaryImage := make(BinaryImage, height)
	for y := 0; y < height; y++ {
		binaryImage[y] = make([]int, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 {
				binaryImage[y][x] = 1
			} else {
				binaryImage[y][x] = 0
			}
		}
	}

	return binaryImage
}

func ConvertIntoImage(binaryImage BinaryImage) *image.RGBA {
	height := len(binaryImage)
	width := len(binaryImage[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if binaryImage[y][x] == 1 {
				img.Set(x, y, image.White)
			} else {
				img.Set(x, y, image.Black)
			}
		}
	}

	return img
}

func Fits(image BinaryImage, se StructuringElement, x, y int) bool {
	rows := len(image)
	cols := len(image[0])

	for i := 0; i < len(se.Data); i++ {
		for j := 0; j < len(se.Data[i]); j++ {
			if se.Data[i][j] == 1 {
				newX := x + i - se.OriginX
				newY := y + j - se.OriginY
				if newX < 0 || newX >= rows || newY < 0 || newY >= cols || image[newX][newY] == 0 {
					return false
				}
			}
		}
	}
	return true
}

func LoadStructureElementsFromJSON(filepath string) (map[string]StructuringElement, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var structureElementsJSON StructureElementsJSON
	err = json.Unmarshal(data, &structureElementsJSON)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return structureElementsJSON.StructureElements, nil
}
