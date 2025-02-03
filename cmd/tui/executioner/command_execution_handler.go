package executioner

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type CommandExecutionHandler func(imgPath string, args map[string]string) (string, error)

var commandRegistry = map[string]CommandExecutionHandler{
	"bandpass": bandpassExecutioner,
	"lowpass":  lowpassExecutioner,
	"highpass": highpassExecutioner,
	"bandcut":  bandcutExecutioner,
	"phasemod": phasemodExecutioner,
	"maskpass": maskpassExecutioner,
}

func parseIntArg(args map[string]string, key string) (int, error) {
	value, exists := args[key]
	if !exists {
		return 0, fmt.Errorf("missing required argument: %s", key)
	}
	return strconv.Atoi(value)
}

func parseBoolArg(args map[string]string, key string) (bool, error) {
	value, exists := args[key]
	if !exists {
		return false, fmt.Errorf("missing required argument: %s", key)
	}
	return strconv.ParseBool(value)
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
		return "", errors.New("maskName cannot be empty")
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
