package executioner

import (
	"errors"
	"fmt"
	"image-processing/cmd"
	"image-processing/imageio"
	"image-processing/manipulations"
	"image-processing/orthogonal_transforms"
)

type handlingCommandOptions struct {
	imgPath                                                 string
	maskPath                                                string
	lowCut, highCut, brightnessPercentage, contrast, factor int
	cutoff, k, l                                            int
	withSpectrumImgGenerated                                bool
}

func validateFrequencyRange(lowCut, highCut, imgWidth int) error {
	if lowCut < 0 || highCut < 0 {
		return errors.New("low and high cut frequencies must be positive")
	}

	if lowCut >= highCut {
		return errors.New("low cut frequency must be lower than high cut frequency")
	}

	if lowCut >= imgWidth/2 || highCut >= imgWidth/2 {
		return errors.New("low and high cut frequencies must be lower than half the image width")
	}

	return nil
}

func saveFilteringResults(handlerResult []cmd.ResultImage) error {
	for _, resultImg := range handlerResult {
		err := imageio.SaveBmpImage(resultImg.GetImage(), resultImg.GetName())
		if err != nil {
			return err
		}
	}

	return nil
}

func handleBrightnessCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	bp := opts.brightnessPercentage

	if bp < -100 || bp > 100 {
		return "", errors.New("brightness percentage must be between -100 and 100")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_brightness_%d_altered.bmp", imgFileName, bp)

	adjustedImg := manipulations.AdjustBrightness(img, bp)

	brightnessResult := cmd.BasicImgResult{
		Img:  adjustedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{brightnessResult}); err != nil {
		return "", err
	}

	return "Brightness adjusted successfully", nil
}

func handleContrastCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	contrast := opts.contrast

	if contrast < -255 || contrast > 255 {
		return "", errors.New("contrast value must be between -255 and 255")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_contrast_%d_altered.bmp", imgFileName, contrast)
	resultImg := manipulations.AdjustContrast(img, contrast)

	contrastResult := cmd.BasicImgResult{
		Img:  resultImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{contrastResult}); err != nil {
		return "", err
	}

	return "Contrast adjusted successfully", nil
}

func handleNegativeCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_negative.bmp", imgFileName)

	negativeResult := manipulations.NegativeImage(img)

	negativeImgResult := cmd.BasicImgResult{
		Img:  negativeResult,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{negativeImgResult}); err != nil {
		return "", err
	}

	return "Negative image created successfully", nil
}

func handleHorizontalFlipCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_horizontal_flipped.bmp", imgFileName)

	flippedImg := manipulations.HorizontalFlip(img)

	flipResult := cmd.BasicImgResult{
		Img:  flippedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{flipResult}); err != nil {
		return "", err
	}

	return "Image flipped horizontally successfully", nil
}

func handleVerticalFlipCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_vertical_flipped.bmp", imgFileName)

	flippedImg := manipulations.VerticalFlip(img)

	flipResult := cmd.BasicImgResult{
		Img:  flippedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{flipResult}); err != nil {
		return "", err
	}

	return "Image flipped vertically successfully", nil
}

func handleDiagonalFlipCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_diagonal_flipped.bmp", imgFileName)

	flippedImg := manipulations.DiagonalFlip(img)

	flipResult := cmd.BasicImgResult{
		Img:  flippedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{flipResult}); err != nil {
		return "", err
	}

	return "Image flipped diagonally successfully", nil
}

func handleShrinkCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	shrinkFactor := opts.factor
	outputFileName := fmt.Sprintf("%s_shrunk_%d_times.bmp", imgFileName, shrinkFactor)

	shrunkImg, err := manipulations.ShrinkImage(img, shrinkFactor)

	if err != nil {
		return "", err
	}

	shrinkResult := cmd.BasicImgResult{
		Img:  shrunkImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{shrinkResult}); err != nil {
		return "", err
	}

	return "Image shrunk successfully", nil
}

func handleEnlargeCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	enlargeFactor := opts.factor
	outputFileName := fmt.Sprintf("%s_enlarged_%d_times.bmp", imgFileName, enlargeFactor)

	enlargedImg, err := manipulations.EnlargeImage(img, enlargeFactor)

	if err != nil {
		return "", err
	}

	enlargeResult := cmd.BasicImgResult{
		Img:  enlargedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{enlargeResult}); err != nil {
		return "", err
	}

	return "Image enlarged successfully", nil
}

func handleBandpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if err := validateFrequencyRange(opts.lowCut, opts.highCut, img.Bounds().Dx()); err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleBandpassFiltering(img, imgFileName, opts.lowCut, opts.highCut, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(cmd.ToResultImage(handlerResult)); err != nil {
		return "", err
	}

	return "Bandpass filtering successful", nil
}

func handleLowpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.cutoff < 0 || opts.cutoff >= img.Bounds().Dx()/2 || opts.cutoff >= img.Bounds().Dy()/2 {
		return "", errors.New("cutoff frequency must be positive and less than half the image width/height")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleLowpassFiltering(img, imgFileName, opts.cutoff, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(cmd.ToResultImage(handlerResult)); err != nil {
		return "", err
	}

	return "Lowpass filtering successful", nil
}

func handleHighpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.cutoff < 0 || opts.cutoff >= img.Bounds().Dx()/2 || opts.cutoff >= img.Bounds().Dy()/2 {
		return "", errors.New("cutoff frequency must be positive and less than half the image width/height")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleHighpassFiltering(img, imgFileName, opts.cutoff, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(cmd.ToResultImage(handlerResult)); err != nil {
		return "", err
	}

	return "Highpass filtering successful", nil
}

func handleBandcutCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if err := validateFrequencyRange(opts.lowCut, opts.highCut, img.Bounds().Dx()); err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleBandcutFiltering(img, imgFileName, opts.lowCut, opts.highCut, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(cmd.ToResultImage(handlerResult)); err != nil {
		return "", err
	}

	return "Bandcut filtering successful", nil
}

func handlePhasemodCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	if opts.k < 0 || opts.l < 0 {
		return "", errors.New("k and l must be positive")
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandlePhaseModification(img, imgFileName, opts.k, opts.l)

	if err := saveFilteringResults(cmd.ToResultImage(handlerResult)); err != nil {
		return "", err
	}

	return "Phase modified successful", nil
}

func handleMaskpassCommand(opts handlingCommandOptions) (successMsgString string, err error) {

	maskImg, err := imageio.OpenBmpImage(opts.maskPath)
	if err != nil {
		return "", fmt.Errorf("failed to open mask image: %w", err)
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)

	handlerResult := orthogonal_transforms.HandleMaskpassFiltering(img, imgFileName, maskImg, opts.withSpectrumImgGenerated)

	if err := saveFilteringResults(cmd.ToResultImage(handlerResult)); err != nil {
		return "", err
	}

	return "Maskpass filter successful applied", nil
}
