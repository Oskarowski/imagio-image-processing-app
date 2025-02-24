// taken from unmaintained repository: https://github.com/qeesung/image2ascii

package ascii_preview

import (
	"bytes"
	"image"
	"image/color"
	"imagio/internal/ascii_preview/ascii"
	"log"

	"errors"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
	"github.com/nfnt/resize"
	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

const (
	charWidthWindows = 0.714
	charWidthOther   = 0.5
)

// NewTerminalAccessor create a new terminal accessor
func NewTerminalAccessor() Terminal {
	return Accessor{}
}

// Terminal get the terminal basic information
type Terminal interface {
	CharWidth() float64
	ScreenSize() (width, height int, err error)
	IsWindows() bool
}

// Accessor implement the Terminal interface and
// fetch the terminal basic information
type Accessor struct {
}

// CharWidth get the terminal char width
func (accessor Accessor) CharWidth() float64 {
	if accessor.IsWindows() {
		return charWidthWindows
	}
	return charWidthOther
}

// IsWindows check if current system is windows
func (accessor Accessor) IsWindows() bool {
	return runtime.GOOS == "windows"
}

// ScreenSize get the terminal screen size
func (accessor Accessor) ScreenSize() (newWidth, newHeight int, err error) {
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return 0, 0,
			errors.New("can not detect the terminal")
	}

	x, _ := terminal.Width()
	y, _ := terminal.Height()

	return int(x), int(y), nil
}

// Options to convert the image to ASCII
type Options struct {
	Ratio           float64
	FixedWidth      int
	FixedHeight     int
	FitScreen       bool // only work on terminal
	StretchedScreen bool // only work on terminal
	Colored         bool // only work on terminal
	Reversed        bool
}

// DefaultOptions for convert image
var DefaultOptions = Options{
	Ratio:           1,
	FixedWidth:      -1,
	FixedHeight:     -1,
	FitScreen:       true,
	Colored:         true,
	Reversed:        false,
	StretchedScreen: false,
}

// NewImageConverter create a new image converter
func NewImageConverter() *ImageConverter {
	return &ImageConverter{
		resizeHandler:  NewResizeHandler(),
		pixelConverter: ascii.NewPixelConverter(),
	}
}

// Converter define the convert image basic operations
type Converter interface {
	Image2ASCIIString(image image.Image, options *Options) string
}

// ImageConverter implement the Convert interface, and responsible
// to image conversion
type ImageConverter struct {
	resizeHandler  ResizeHandler
	pixelConverter ascii.PixelConverter
}

// Image2ASCIIMatrix converts a image to ASCII matrix
func (converter *ImageConverter) Image2ASCIIMatrix(image image.Image, imageConvertOptions *Options) []string {
	// Resize the convert first
	newImage := converter.resizeHandler.ScaleImage(image, imageConvertOptions)
	sz := newImage.Bounds()
	newWidth := sz.Max.X
	newHeight := sz.Max.Y
	rawCharValues := make([]string, 0, int(newWidth*newHeight+newWidth))
	for i := 0; i < int(newHeight); i++ {
		for j := 0; j < int(newWidth); j++ {
			pixel := color.NRGBAModel.Convert(newImage.At(j, i))
			// Convert the pixel to ascii char
			pixelConvertOptions := ascii.NewOptions()
			pixelConvertOptions.Colored = imageConvertOptions.Colored
			pixelConvertOptions.Reversed = imageConvertOptions.Reversed
			rawChar := converter.pixelConverter.ConvertPixelToASCII(pixel, &pixelConvertOptions)
			rawCharValues = append(rawCharValues, rawChar)
		}
		rawCharValues = append(rawCharValues, "\n")
	}
	return rawCharValues
}

// Image2ASCIIString converts a image to ascii matrix, and the join the matrix to a string
func (converter *ImageConverter) Image2ASCIIString(image image.Image, options *Options) string {
	convertedPixelASCII := converter.Image2ASCIIMatrix(image, options)
	var buffer bytes.Buffer

	for i := 0; i < len(convertedPixelASCII); i++ {
		buffer.WriteString(convertedPixelASCII[i])
	}
	return buffer.String()
}

// NewResizeHandler create a new resize handler
func NewResizeHandler() ResizeHandler {
	handler := &ImageResizeHandler{
		terminal: NewTerminalAccessor(),
	}

	initResizeResolver(handler)
	return handler
}

// initResizeResolver register the size resolvers
func initResizeResolver(handler *ImageResizeHandler) {
	sizeResolvers := make([]imageSizeResolver, 0, 5)
	// fixed height or width resolver
	sizeResolvers = append(sizeResolvers, imageSizeResolver{
		match: func(options *Options) bool {
			return options.FixedWidth != -1 || options.FixedHeight != -1
		},
		compute: func(sz image.Rectangle, options *Options, handler *ImageResizeHandler) (width, height int, err error) {
			height = sz.Max.Y
			width = sz.Max.X
			if options.FixedWidth != -1 {
				width = options.FixedWidth
			}

			if options.FixedHeight != -1 {
				height = options.FixedHeight
			}
			return
		},
	})
	// scaled by ratio
	sizeResolvers = append(sizeResolvers, imageSizeResolver{
		match: func(options *Options) bool {
			return options.Ratio != 1
		},
		compute: func(sz image.Rectangle, options *Options, handler *ImageResizeHandler) (width, height int, err error) {
			ratio := options.Ratio
			width = handler.ScaleWidthByRatio(float64(sz.Max.X), ratio)
			height = handler.ScaleHeightByRatio(float64(sz.Max.Y), ratio)
			return
		},
	})
	// scaled by stretched screen
	sizeResolvers = append(sizeResolvers, imageSizeResolver{
		match: func(options *Options) bool {
			return options.StretchedScreen
		},
		compute: func(sz image.Rectangle, options *Options, handler *ImageResizeHandler) (width, height int, err error) {
			return handler.terminal.ScreenSize()
		},
	})
	// scaled by fit the screen
	sizeResolvers = append(sizeResolvers, imageSizeResolver{
		match: func(options *Options) bool {
			return options.FitScreen
		},
		compute: func(sz image.Rectangle, options *Options, handler *ImageResizeHandler) (width, height int, err error) {
			return handler.CalcProportionalFittingScreenSize(sz)
		},
	})
	// default size
	sizeResolvers = append(sizeResolvers, imageSizeResolver{
		match: func(options *Options) bool {
			return true
		},
		compute: func(sz image.Rectangle, options *Options, handler *ImageResizeHandler) (width, height int, err error) {
			return sz.Max.X, sz.Max.Y, nil
		},
	})

	handler.imageSizeResolvers = sizeResolvers
}

// ResizeHandler define the operation to resize a image
type ResizeHandler interface {
	ScaleImage(image image.Image, options *Options) (newImage image.Image)
}

// ImageResizeHandler implement the ResizeHandler interface and
// responsible for image resizing
type ImageResizeHandler struct {
	terminal           Terminal
	imageSizeResolvers []imageSizeResolver
}

// imageSizeResolver to resolve the image size from the options
type imageSizeResolver struct {
	match   func(options *Options) bool
	compute func(sz image.Rectangle, options *Options, handler *ImageResizeHandler) (width, height int, err error)
}

// ScaleImage resize the convert to expected size base on the convert options
func (handler *ImageResizeHandler) ScaleImage(image image.Image, options *Options) (newImage image.Image) {
	sz := image.Bounds()
	newWidth, newHeight, err := handler.resolveSize(sz, options)
	if err != nil {
		log.Fatal(err)
	}

	newImage = resize.Resize(uint(newWidth), uint(newHeight), image, resize.Lanczos3)
	return
}

// resolveSize resolve the image size
func (handler *ImageResizeHandler) resolveSize(sz image.Rectangle, options *Options) (width, height int, err error) {
	for _, resolver := range handler.imageSizeResolvers {
		if resolver.match(options) {
			return resolver.compute(sz, options, handler)
		}
	}
	return sz.Max.X, sz.Max.Y, nil
}

// CalcProportionalFittingScreenSize proportional scale the image
// so that the terminal can just show the picture.
func (handler *ImageResizeHandler) CalcProportionalFittingScreenSize(sz image.Rectangle) (newWidth, newHeight int, err error) {
	screenWidth, screenHeight, err := handler.terminal.ScreenSize()
	if err != nil {
		log.Fatal(nil)
	}
	newWidth, newHeight = handler.CalcFitSize(
		float64(screenWidth),
		float64(screenHeight),
		float64(sz.Max.X),
		float64(sz.Max.Y))
	return
}

// CalcFitSizeRatio through the given length and width,
// the computation can match the optimal scaling ratio of the length and width.
// In other words, it is able to give a given size rectangle to contain pictures
// Either match the width first, then scale the length equally,
// or match the length first, then scale the height equally.
// More detail please check the example
func (handler *ImageResizeHandler) CalcFitSizeRatio(width, height, imageWidth, imageHeight float64) (ratio float64) {
	ratio = 1.0
	// try to fit the height
	ratio = height / imageHeight
	scaledWidth := imageWidth * ratio / handler.terminal.CharWidth()
	if scaledWidth < width {
		return ratio / handler.terminal.CharWidth()
	}

	// try to fit the width
	ratio = width / imageWidth
	return ratio
}

// CalcFitSize through the given length and width ,
// Calculation is able to match the length and width of
// the specified size, and is proportional scaling.
func (handler *ImageResizeHandler) CalcFitSize(width, height, toBeFitWidth, toBeFitHeight float64) (fitWidth, fitHeight int) {
	ratio := handler.CalcFitSizeRatio(width, height, toBeFitWidth, toBeFitHeight)
	fitWidth = handler.ScaleWidthByRatio(toBeFitWidth, ratio)
	fitHeight = handler.ScaleHeightByRatio(toBeFitHeight, ratio)
	return
}

// ScaleWidthByRatio scaled the width by ratio
func (handler *ImageResizeHandler) ScaleWidthByRatio(width float64, ratio float64) int {
	return int(width * ratio)
}

// ScaleHeightByRatio scaled the height by ratio
func (handler *ImageResizeHandler) ScaleHeightByRatio(height float64, ratio float64) int {
	return int(height * ratio * handler.terminal.CharWidth())
}
