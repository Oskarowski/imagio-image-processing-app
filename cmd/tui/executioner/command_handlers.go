package executioner

import (
	"errors"
	"fmt"
	"image-processing/analysis"
	"image-processing/cmd"
	"image-processing/imageio"
	"image-processing/manipulations"
	"image-processing/noise"
	"image-processing/orthogonal_transforms"
	"strings"
)

type handlingCommandOptions struct {
	imgPath, maskPath, selectedComparisonCommands, comparisonImagePath string
	lowCut, highCut, brightnessPercentage, contrast, factor            int
	cutoff, k, l, maxWindowSize, minWindowSize                         int
	withSpectrumImgGenerated                                           bool
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

func handleAdaptiveNoiseFilterCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	if opts.minWindowSize < 1 || opts.maxWindowSize < 1 {
		return "", errors.New("window sizes must be greater than 1")
	}

	if opts.minWindowSize > opts.maxWindowSize {
		return "", errors.New("minWindowSize must be less than maxWindowSize")
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.minWindowSize > img.Bounds().Dx()/2 || opts.maxWindowSize > img.Bounds().Dx()/2 {
		return "", errors.New("window sizes must be less than half the image width")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_denoised_via_adaptive_filter.bmp", imgFileName)

	denoisedImg := noise.AdaptiveMedianFilter(img, opts.minWindowSize, opts.maxWindowSize)

	denoisedResult := cmd.BasicImgResult{
		Img:  denoisedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{denoisedResult}); err != nil {
		return "", err
	}

	return "Image denoised via adaptive filter successfully", nil
}

func handleMinNoiseFilterCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	if opts.minWindowSize < 1 {
		return "", errors.New("window size must be greater than 1")
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.minWindowSize > img.Bounds().Dx()/2 {
		return "", errors.New("window size must be less than half the image width")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_denoised_via_min_filter.bmp", imgFileName)

	denoisedImg := noise.MinFilter(img, opts.minWindowSize)

	denoisedResult := cmd.BasicImgResult{
		Img:  denoisedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{denoisedResult}); err != nil {
		return "", err
	}

	return "Image denoised via min filter successfully", nil
}

func handleMaxNoiseFilterCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	if opts.maxWindowSize < 1 {
		return "", errors.New("window size must be greater than 1")
	}

	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	if opts.maxWindowSize > img.Bounds().Dx()/2 {
		return "", errors.New("window size must be less than half the image width")
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_denoised_via_max_filter.bmp", imgFileName)

	denoisedImg := noise.MaxFilter(img, opts.maxWindowSize)

	denoisedResult := cmd.BasicImgResult{
		Img:  denoisedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{denoisedResult}); err != nil {
		return "", err
	}

	return "Image denoised via max filter successfully", nil
}

func handleImgComparisonCommand(opts handlingCommandOptions) (successMsgString string, output []analysis.ImageComparisonEntry, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", nil, err
	}

	comparisonImg, err := imageio.OpenBmpImage(opts.comparisonImagePath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open comparison image: %w", err)
	}

	selectedCommands := strings.Split(opts.selectedComparisonCommands, "|")

	if len(selectedCommands) == 0 {
		return "", nil, errors.New("no comparison methods selected")
	}

	var entries []analysis.ImageComparisonEntry

	img1Name := imageio.GetFileName(opts.imgPath)
	img2Name := imageio.GetFileName(opts.comparisonImagePath)

	for _, comparison := range selectedCommands {
		result := analysis.CalculateComparisonCharacteristic(comparison, img, comparisonImg)
		result.Img1Name = img1Name
		result.Img2Name = img2Name
		entries = append(entries, result)
	}

	return "Images compared successfully", entries, nil
}

func handleImgHistogramCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	histogramImg := manipulations.GenerateGraphicalRepresentationOfHistogram(manipulations.CalculateHistogram(img))

	histogramResult := cmd.BasicImgResult{
		Img:  histogramImg,
		Name: fmt.Sprintf("%s_histogram.bmp", imgFileName),
	}

	if err := saveFilteringResults([]cmd.ResultImage{histogramResult}); err != nil {
		return "", err
	}

	return "Histogram calculated successfully", nil
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
