package manipulations

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	paddingLeft   = 60
	paddingRight  = 20
	paddingBottom = 30
	paddingTop    = 20
	yLabelOffset  = 10
	isXAxis       = true
	isYAxis       = false
)

func CalculateHistogram(img image.Image, channel int) *image.RGBA {
	var histogram [256]int
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			var intensity uint8
			switch channel {
			case 0:
				intensity = pixel.R
			case 1:
				intensity = pixel.G
			case 2:
				intensity = pixel.B
			default:
				return nil
			}
			histogram[intensity]++

		}
	}

	maxFrequency := 0
	for _, freq := range histogram {
		if freq > maxFrequency {
			maxFrequency = freq
		}
	}

	width, height := 500, 500
	histImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(histImg, histImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw histogram bars with scaling
	barColor := color.RGBA{50, 100, 245, 255}
	for x := 0; x < 256; x++ {
		frequency := histogram[x]
		scaleHeight := int(float64(frequency) / float64(maxFrequency) * float64(height-paddingBottom-paddingTop))

		for y := height - paddingBottom; y >= height-scaleHeight-paddingBottom; y-- {
			histImg.Set(x+paddingLeft, y, barColor)
		}
	}

	drawAxes(histImg, maxFrequency)
	return histImg
}

func drawAxes(img *image.RGBA, maxY int) {
	black := color.RGBA{0, 0, 0, 255}
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	// Draw x-axis (intensity) and y-axis (frequency)
	for x := paddingLeft; x < width-paddingRight; x++ {
		img.Set(x, height-paddingBottom, black)
	}
	for y := paddingTop; y < height-paddingBottom; y++ {
		img.Set(paddingLeft, y, black)
	}

	// Add tick marks and labels on x-axis (0-255 intensity values)
	for x := 0; x <= 255; x += 50 {
		for y := height - paddingBottom - 5; y < height-paddingBottom+5; y++ {
			img.Set(x+paddingLeft, y, black)
		}
		labelValue(img, x, height-paddingBottom+15, black, isXAxis)
	}

	// Add tick marks and labels on y-axis based on max frequency scaling
	yTicks := 5
	yStep := maxY / yTicks
	yLabelOffset := basicfont.Face7x13.Ascent / 2
	for i := 0; i <= yTicks; i++ {
		yPos := height - paddingBottom - (i * (height - paddingBottom - paddingTop) / yTicks)
		for x := paddingLeft - 5; x < paddingLeft+5; x++ {
			img.Set(x, yPos, black)
		}
		labelValue(img, i*yStep, yPos+yLabelOffset, black, isYAxis)
	}
}

func labelValue(img *image.RGBA, value int, y int, col color.Color, isXAxis bool) {
	str := fmt.Sprintf("%d", value)

	// Position text based on axis
	var x int
	if isXAxis {
		x = value + paddingLeft - (len(str) * basicfont.Face7x13.Width / 2)
	} else {
		x = paddingLeft - yLabelOffset - (len(str) * basicfont.Face7x13.Width)
	}

	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	drawer := &font.Drawer{
		Dst:  img,
		Src:  &image.Uniform{col},
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	drawer.DrawString(str)
}
