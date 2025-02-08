package executioner

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type CommandExecutionHandler func(imgPath string, args map[string]string) (string, error)

var commandRegistry = map[string]CommandExecutionHandler{
	"brightness":        brightnessExecutioner,
	"contrast":          contrastExecutioner,
	"negative":          negativeExecutioner,
	"flip_horizontally": flipHorizontallyExecutioner,
	"flip_vertically":   flipVerticallyExecutioner,
	"flip_diagonally":   flipDiagonallyExecutioner,
	"bandpass":          bandpassExecutioner,
	"lowpass":           lowpassExecutioner,
	"highpass":          highpassExecutioner,
	"bandcut":           bandcutExecutioner,
	"phasemod":          phasemodExecutioner,
	"maskpass":          maskpassExecutioner,
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
