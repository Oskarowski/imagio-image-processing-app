package executioner

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type CommandExecutionHandler func(imgPath string, args map[string]string) (string, error)

var commandRegistry = map[string]CommandExecutionHandler{
	"brightness":                brightnessExecutioner,
	"contrast":                  contrastExecutioner,
	"negative":                  negativeExecutioner,
	"flip_horizontally":         flipHorizontallyExecutioner,
	"flip_vertically":           flipVerticallyExecutioner,
	"flip_diagonally":           flipDiagonallyExecutioner,
	"shrink":                    shrinkExecutioner,
	"enlarge":                   enlargeExecutioner,
	"adaptive_filter_denoising": adaptiveNoiseFilterExecutioner,
	"min_filter_denoising":      minNoiseFilterExecutioner,
	"max_filter_denoising":      maxNoiseFilterExecutioner,
	"img_comparison_commands":   imgComparisonExecutioner,
	"bandpass":                  bandpassExecutioner,
	"lowpass":                   lowpassExecutioner,
	"highpass":                  highpassExecutioner,
	"bandcut":                   bandcutExecutioner,
	"phasemod":                  phasemodExecutioner,
	"maskpass":                  maskpassExecutioner,
}

func brightnessExecutioner(imgPath string, args map[string]string) (string, error) {
	brightness, err := parseIntArg(args, "brightness")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:              imgPath,
		brightnessPercentage: brightness,
	}

	return handleBrightnessCommand(opts)
}

func contrastExecutioner(imgPath string, args map[string]string) (string, error) {
	contrast, err := parseIntArg(args, "contrast")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:  imgPath,
		contrast: contrast,
	}

	return handleContrastCommand(opts)
}

func negativeExecutioner(imgPath string, args map[string]string) (string, error) {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	return handleNegativeCommand(opts)
}

func flipHorizontallyExecutioner(imgPath string, args map[string]string) (string, error) {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	return handleHorizontalFlipCommand(opts)
}

func flipVerticallyExecutioner(imgPath string, args map[string]string) (string, error) {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	return handleVerticalFlipCommand(opts)
}

func flipDiagonallyExecutioner(imgPath string, args map[string]string) (string, error) {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	return handleDiagonalFlipCommand(opts)
}

func shrinkExecutioner(imgPath string, args map[string]string) (string, error) {
	factor, err := parseIntArg(args, "shrinkFactor")
	if err != nil {
		return "", err
	}

	if factor < 1 {
		return "", errors.New("shrink factor must be greater than 1")
	}

	opts := handlingCommandOptions{
		imgPath: imgPath,
		factor:  factor,
	}

	return handleShrinkCommand(opts)
}

func enlargeExecutioner(imgPath string, args map[string]string) (string, error) {
	factor, err := parseIntArg(args, "enlargeFactor")
	if err != nil {
		return "", err
	}

	if factor < 1 {
		return "", errors.New("enlarge factor must be greater than 1")
	}

	opts := handlingCommandOptions{
		imgPath: imgPath,
		factor:  factor,
	}

	return handleEnlargeCommand(opts)
}

func adaptiveNoiseFilterExecutioner(imgPath string, args map[string]string) (string, error) {
	minWindowSize, err := parseIntArg(args, "minWindowSize")
	if err != nil {
		return "", err
	}

	maxWindowSize, err := parseIntArg(args, "maxWindowSize")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:       imgPath,
		minWindowSize: minWindowSize,
		maxWindowSize: maxWindowSize,
	}

	return handleAdaptiveNoiseFilterCommand(opts)
}

func minNoiseFilterExecutioner(imgPath string, args map[string]string) (string, error) {
	minWindowSize, err := parseIntArg(args, "minWindowSize")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:       imgPath,
		minWindowSize: minWindowSize,
	}

	return handleMinNoiseFilterCommand(opts)
}

func maxNoiseFilterExecutioner(imgPath string, args map[string]string) (string, error) {
	maxWindowSize, err := parseIntArg(args, "maxWindowSize")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:       imgPath,
		maxWindowSize: maxWindowSize,
	}

	return handleMaxNoiseFilterCommand(opts)
}

func imgComparisonExecutioner(imgPath string, args map[string]string) (string, error) {
	opts := handlingCommandOptions{
		imgPath:                    imgPath,
		comparisonImagePath:        args["comparisonImagePath"],
		selectedComparisonCommands: args["selectedComparisonCommands"],
	}

	return handleImgComparisonCommand(opts)
}

func bandpassExecutioner(imgPath string, args map[string]string) (string, error) {
	lowCut, err := parseIntArg(args, "lowCut")
	if err != nil {
		return "", err
	}

	highCut, err := parseIntArg(args, "highCut")
	if err != nil {
		return "", err
	}

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		lowCut:                   lowCut,
		highCut:                  highCut,
		withSpectrumImgGenerated: withSpectrum,
	}

	return handleBandpassCommand(opts)
}

func lowpassExecutioner(imgPath string, args map[string]string) (string, error) {
	cutoff, err := parseIntArg(args, "cutoff")
	if err != nil {
		return "", err
	}

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		cutoff:                   cutoff,
		withSpectrumImgGenerated: withSpectrum,
	}

	return handleLowpassCommand(opts)
}

func highpassExecutioner(imgPath string, args map[string]string) (string, error) {
	cutoff, err := parseIntArg(args, "cutoff")
	if err != nil {
		return "", err
	}

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		cutoff:                   cutoff,
		withSpectrumImgGenerated: withSpectrum,
	}

	return handleHighpassCommand(opts)
}

func bandcutExecutioner(imgPath string, args map[string]string) (string, error) {
	lowCut, err := parseIntArg(args, "lowCut")
	if err != nil {
		return "", err
	}
	highCut, err := parseIntArg(args, "highCut")
	if err != nil {
		return "", err
	}
	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		lowCut:                   lowCut,
		highCut:                  highCut,
		withSpectrumImgGenerated: withSpectrum,
	}

	return handleBandcutCommand(opts)
}

func phasemodExecutioner(imgPath string, args map[string]string) (string, error) {
	k, err := parseIntArg(args, "k")
	if err != nil {
		return "", err
	}

	l, err := parseIntArg(args, "l")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath: imgPath,
		k:       k,
		l:       l,
	}

	return handlePhasemodCommand(opts)
}

func maskpassExecutioner(imgPath string, args map[string]string) (string, error) {
	maskName := strings.TrimSpace(args["maskName"])
	if maskName == "" {
		return "", errors.New("mask name cannot be empty")
	}

	maskPath := filepath.Join("orthogonal_transforms", "masks", maskName)

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return "", err
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		maskPath:                 maskPath,
		withSpectrumImgGenerated: withSpectrum,
	}

	return handleMaskpassCommand(opts)
}

func ExecuteCommand(imgPath, cmdName string, cmdArgs map[string]string) (string, error) {
	handler, exists := commandRegistry[cmdName]
	if !exists {
		return "", fmt.Errorf("command not found: %s", cmdName)
	}
	return handler(imgPath, cmdArgs)
}
