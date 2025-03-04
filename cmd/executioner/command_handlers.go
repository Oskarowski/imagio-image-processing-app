package executioner

import (
	"errors"
	"fmt"
	"imagio/analysis"
	"imagio/cmd"
	"imagio/imageio"
	"imagio/manipulations"
	"imagio/morphological"
	"imagio/noise"
	"imagio/orthogonal_transforms"
	"strings"
)

type handlingCommandOptions struct {
	imgPath, maskPath, selectedComparisonCommands, comparisonImagePath, selectedHistogramCharacteristicsCommands, maskName, structureElementName, foregroundSE, backgroundSE, seedPointsStr string
	lowCut, highCut, brightnessPercentage, contrast, factor, distanceMetric                                                                                                                 int
	cutoff, k, l, maxWindowSize, minWindowSize                                                                                                                                              int
	alphaValue, thresholdValue                                                                                                                                                              float64
	withSpectrumImgGenerated                                                                                                                                                                bool
	edgeSharpeningMask                                                                                                                                                                      [][]int
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

func handleImgComparisonCommand(opts handlingCommandOptions) (successMsgString string, output []analysis.CharacteristicsEntry, err error) {
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

	var entries []analysis.CharacteristicsEntry

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

func handleHistogramImgCharacteristicsCommand(opts handlingCommandOptions) (successMsgString string, output []analysis.CharacteristicsEntry, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", nil, err
	}

	selectedCharacteristics := strings.Split(opts.selectedHistogramCharacteristicsCommands, "|")

	if len(selectedCharacteristics) == 0 {
		return "", nil, errors.New("no histogram characteristics selected")
	}

	imgPureName := imageio.GetPureFileName(opts.imgPath)
	imgName := imageio.GetFileName(opts.imgPath)

	var characteristics []analysis.CharacteristicsEntry

	for _, characteristic := range selectedCharacteristics {
		histogram := manipulations.CalculateHistogram(img)
		result := analysis.CalculateHistogramCharacteristic(characteristic, histogram, imgPureName)
		result.Img1Name = imgName
		characteristics = append(characteristics, result)
	}

	return "Histogram characteristics calculated successfully", characteristics, nil
}

func handleRayleighTransformCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	gMin := opts.lowCut
	gMax := opts.highCut
	alpha := opts.alphaValue

	if gMin < 0 || gMax > 255 || gMin >= gMax {
		return "", errors.New("lowCutoff and highCutoff must be between 0 and 255 and lowCutoff must be less than highCutoff")
	}

	rayleighTransformedImg := manipulations.EnhanceImageWithRayleigh(img, float64(gMin), float64(gMax), alpha)

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_rayleigh_transformed_min%d_max%d_alpha%.2f.bmp", imgFileName, gMin, gMax, alpha)

	rayleighResult := cmd.BasicImgResult{
		Img:  rayleighTransformedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{rayleighResult}); err != nil {
		return "", err
	}

	successMsg := fmt.Sprintf("Rayleigh transformation applied successfully with min: %d, max: %d, alpha: %.2f", gMin, gMax, alpha)

	return successMsg, nil
}

func handleMaskEdgeSharpeningCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	sharpenedImg := manipulations.ApplyConvolutionUniversal(img, opts.edgeSharpeningMask)

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_sharpened_with_%s_mask.bmp", imgFileName, opts.maskName)

	sharpenedResult := cmd.BasicImgResult{
		Img:  sharpenedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{sharpenedResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Edge sharpening mask applied successfully with mask: %s", opts.maskName)
	return msg, nil
}

func handleKirshEdgeDetectionCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	edgeDetectedImg := manipulations.ApplyKirshEdgeDetection(img)

	imgFileName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_kirsh_edge_detected.bmp", imgFileName)

	edgeDetectedResult := cmd.BasicImgResult{
		Img:  edgeDetectedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{edgeDetectedResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Kirsh edge detection applied successfully to %s", imgFileName)
	return msg, nil
}

func handleDilationCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	structuringElement, err := morphological.GetStructureElement(opts.structureElementName)
	if err != nil {
		return "", err
	}

	binaryImg := morphological.ConvertIntoBinaryImage(img)
	dilatedBinaryImg := morphological.Dilation(binaryImg, structuringElement)
	dilatedImg := morphological.ConvertIntoImage(dilatedBinaryImg)

	pureImgName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_dilated_with_%s.bmp", pureImgName, opts.structureElementName)

	dilatedResult := cmd.BasicImgResult{
		Img:  dilatedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{dilatedResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Dilation applied successfully with %s structural element", opts.structureElementName)
	return msg, nil
}

func handleErosionCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	structuringElement, err := morphological.GetStructureElement(opts.structureElementName)
	if err != nil {
		return "", err
	}

	binaryImg := morphological.ConvertIntoBinaryImage(img)
	dilatedBinaryImg := morphological.Erosion(binaryImg, structuringElement)
	dilatedImg := morphological.ConvertIntoImage(dilatedBinaryImg)

	pureImgName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_eroded_with_%s.bmp", pureImgName, opts.structureElementName)

	dilatedResult := cmd.BasicImgResult{
		Img:  dilatedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{dilatedResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Erosion applied successfully with %s structural element", opts.structureElementName)
	return msg, nil
}

func handleOpeningCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	structuringElement, err := morphological.GetStructureElement(opts.structureElementName)
	if err != nil {
		return "", err
	}

	binaryImg := morphological.ConvertIntoBinaryImage(img)
	dilatedBinaryImg := morphological.Opening(binaryImg, structuringElement)
	dilatedImg := morphological.ConvertIntoImage(dilatedBinaryImg)

	pureImgName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_opened_with_%s.bmp", pureImgName, opts.structureElementName)

	dilatedResult := cmd.BasicImgResult{
		Img:  dilatedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{dilatedResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Opening applied successfully with %s structural element", opts.structureElementName)
	return msg, nil
}

func handleClosingCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	structuringElement, err := morphological.GetStructureElement(opts.structureElementName)
	if err != nil {
		return "", err
	}

	binaryImg := morphological.ConvertIntoBinaryImage(img)
	dilatedBinaryImg := morphological.Closing(binaryImg, structuringElement)
	dilatedImg := morphological.ConvertIntoImage(dilatedBinaryImg)

	pureImgName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_closed_with_%s.bmp", pureImgName, opts.structureElementName)

	dilatedResult := cmd.BasicImgResult{
		Img:  dilatedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{dilatedResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Closing applied successfully with %s structural element", opts.structureElementName)
	return msg, nil
}

func handleHitOrMissCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	fse, err := morphological.GetStructureElement(opts.foregroundSE)
	if err != nil {
		return "", err
	}

	bse, err := morphological.GetStructureElement(opts.backgroundSE)
	if err != nil {
		return "", err
	}

	binaryImg := morphological.ConvertIntoBinaryImage(img)
	hitOrMissBinaryImg := morphological.HitOrMiss(binaryImg, fse, bse)
	hitOrMissImg := morphological.ConvertIntoImage(hitOrMissBinaryImg)

	pureImgName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_hit_or_miss_fse_%s_bse_%s.bmp", pureImgName, opts.foregroundSE, opts.backgroundSE)

	hitOrMissResult := cmd.BasicImgResult{
		Img:  hitOrMissImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{hitOrMissResult}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Hit or miss applied successfully with foreground SE: %s and background SE: %s", opts.foregroundSE, opts.backgroundSE)
	return msg, nil
}

func handleThinningCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	binaryImg := morphological.ConvertIntoBinaryImage(img)
	seSeries := morphological.SeriesXIISE
	thinnedBinaryImg := morphological.Thinning(binaryImg, seSeries)
	thinnedImg := morphological.ConvertIntoImage(thinnedBinaryImg)

	pureImgName := imageio.GetPureFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_thinned.bmp", pureImgName)

	thinnedResult := cmd.BasicImgResult{
		Img:  thinnedImg,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{thinnedResult}); err != nil {
		return "", err
	}

	return "Thinning applied successfully", nil
}

func handleRegionGrowCommand(opts handlingCommandOptions) (successMsgString string, err error) {
	img, err := imageio.OpenBmpImage(opts.imgPath)
	if err != nil {
		return "", err
	}

	seeds, err := morphological.ParseSeedPoints(opts.seedPointsStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse seed points: %w", err)
	}

	metric := morphological.DistanceCriterion(opts.distanceMetric)
	threshold := opts.thresholdValue

	if metric < 0 || metric > 2 {
		return "", errors.New("metric must be between 0 and 2")
	}

	if threshold < 0 {
		return "", errors.New("threshold must be positive number")
	}

	_, imgWithMarkedRegions := morphological.RegionGrowing(img, seeds, metric, threshold)

	imgName := imageio.GetFileName(opts.imgPath)
	outputFileName := fmt.Sprintf("%s_region_grown.bmp", imgName)

	result := cmd.BasicImgResult{
		Img:  imgWithMarkedRegions,
		Name: outputFileName,
	}

	if err := saveFilteringResults([]cmd.ResultImage{result}); err != nil {
		return "", err
	}

	msg := fmt.Sprintf("Region growing applied successfully with metric: %d and threshold: %.2f", metric, threshold)
	return msg, nil
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
