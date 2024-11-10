package manipulations

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

func LoadMasksFromJSON(filename string) (map[string][][]int, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read JSON file: %w", err)
	}

	// TODO add validation for masks format
	var masks map[string][][]int
	if err := json.Unmarshal(bytes, &masks); err != nil {
		return nil, fmt.Errorf("could not parse JSON: %w", err)
	}
	return masks, nil
}

func GetMask(masks map[string][][]int, maskName string) ([][]int, error) {
	if mask, exists := masks[maskName]; exists {
		return mask, nil
	}
	return nil, fmt.Errorf("mask %s not found", maskName)
}

func ApplyConvolutionUniversal(img image.Image, mask [][]int) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	offset := len(mask) / 2

	grayImg := image.NewGray(bounds)
	draw.Draw(grayImg, bounds, img, bounds.Min, draw.Src)

	for y := offset; y < height-offset; y++ {
		for x := offset; x < width-offset; x++ {
			var sum int

			for i := -offset; i <= offset; i++ {
				for j := -offset; j <= offset; j++ {
					grayPixel := color.GrayModel.Convert(img.At(x+i, y+j)).(color.Gray)
					sum += int(grayPixel.Y) * mask[offset+i][offset+j]
				}
			}

			grayImg.SetGray(x, y, color.Gray{Y: ClampUint8(sum)})
		}
	}

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, grayImg, bounds.Min, draw.Src)

	return rgbaImg

}

func ApplyConvolutionOptimized(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	offset := 1

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)

	for y := offset; y < height-offset; y++ {
		for x := offset; x < width-offset; x++ {

			centerGray := color.GrayModel.Convert(img.At(x, y)).(color.Gray).Y
			sum := int(centerGray) * 9

			sum -= int(color.GrayModel.Convert(img.At(x-1, y)).(color.Gray).Y)   // left
			sum -= int(color.GrayModel.Convert(img.At(x+1, y)).(color.Gray).Y)   // right
			sum -= int(color.GrayModel.Convert(img.At(x, y-1)).(color.Gray).Y)   // top
			sum -= int(color.GrayModel.Convert(img.At(x, y+1)).(color.Gray).Y)   // bottom
			sum -= int(color.GrayModel.Convert(img.At(x-1, y-1)).(color.Gray).Y) // top-left
			sum -= int(color.GrayModel.Convert(img.At(x+1, y-1)).(color.Gray).Y) // top-right
			sum -= int(color.GrayModel.Convert(img.At(x-1, y+1)).(color.Gray).Y) // bottom-left
			sum -= int(color.GrayModel.Convert(img.At(x+1, y+1)).(color.Gray).Y) // bottom-right

			clampedValue := ClampUint8(sum)
			rgbaImg.Set(x, y, color.RGBA{clampedValue, clampedValue, clampedValue, 255})
		}
	}

	return rgbaImg
}
