package executioner

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type ExecutionResult struct {
	Message string
	Err     error
	Output  interface{}
}

type CommandExecutionHandler func(imgPath string, args map[string]string) ExecutionResult

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

func brightnessExecutioner(imgPath string, args map[string]string) ExecutionResult {
	brightness, err := parseIntArg(args, "brightness")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:              imgPath,
		brightnessPercentage: brightness,
	}

	msg, err := handleBrightnessCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}

}

func contrastExecutioner(imgPath string, args map[string]string) ExecutionResult {
	contrast, err := parseIntArg(args, "contrast")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:  imgPath,
		contrast: contrast,
	}

	msg, err := handleContrastCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func negativeExecutioner(imgPath string, args map[string]string) ExecutionResult {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	msg, err := handleNegativeCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func flipHorizontallyExecutioner(imgPath string, args map[string]string) ExecutionResult {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	msg, err := handleHorizontalFlipCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func flipVerticallyExecutioner(imgPath string, args map[string]string) ExecutionResult {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	msg, err := handleVerticalFlipCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func flipDiagonallyExecutioner(imgPath string, args map[string]string) ExecutionResult {
	opts := handlingCommandOptions{
		imgPath: imgPath,
	}

	msg, err := handleDiagonalFlipCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func shrinkExecutioner(imgPath string, args map[string]string) ExecutionResult {
	factor, err := parseIntArg(args, "shrinkFactor")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	if factor < 1 {
		return ExecutionResult{
			Message: "",
			Err:     errors.New("shrink factor must be greater than 1"),
		}
	}

	opts := handlingCommandOptions{
		imgPath: imgPath,
		factor:  factor,
	}

	msg, err := handleShrinkCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func enlargeExecutioner(imgPath string, args map[string]string) ExecutionResult {
	factor, err := parseIntArg(args, "enlargeFactor")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	if factor < 1 {
		return ExecutionResult{
			Message: "",
			Err:     errors.New("enlarge factor must be greater than 1"),
		}
	}

	opts := handlingCommandOptions{
		imgPath: imgPath,
		factor:  factor,
	}

	msg, err := handleEnlargeCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func adaptiveNoiseFilterExecutioner(imgPath string, args map[string]string) ExecutionResult {
	minWindowSize, err := parseIntArg(args, "minWindowSize")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	maxWindowSize, err := parseIntArg(args, "maxWindowSize")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:       imgPath,
		minWindowSize: minWindowSize,
		maxWindowSize: maxWindowSize,
	}

	msg, err := handleAdaptiveNoiseFilterCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func minNoiseFilterExecutioner(imgPath string, args map[string]string) ExecutionResult {
	minWindowSize, err := parseIntArg(args, "minWindowSize")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:       imgPath,
		minWindowSize: minWindowSize,
	}

	msg, err := handleMinNoiseFilterCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func maxNoiseFilterExecutioner(imgPath string, args map[string]string) ExecutionResult {
	maxWindowSize, err := parseIntArg(args, "maxWindowSize")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:       imgPath,
		maxWindowSize: maxWindowSize,
	}

	msg, err := handleMaxNoiseFilterCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func imgComparisonExecutioner(imgPath string, args map[string]string) ExecutionResult {
	opts := handlingCommandOptions{
		imgPath:                    imgPath,
		comparisonImagePath:        args["comparisonImagePath"],
		selectedComparisonCommands: args["selectedComparisonCommands"],
	}

	log.Default().Println(opts.selectedComparisonCommands)

	msg, output, err := handleImgComparisonCommand(opts)

	return ExecutionResult{
		Message: msg,
		Output:  output,
		Err:     err,
	}
}

func bandpassExecutioner(imgPath string, args map[string]string) ExecutionResult {
	lowCut, err := parseIntArg(args, "lowCut")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	highCut, err := parseIntArg(args, "highCut")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		lowCut:                   lowCut,
		highCut:                  highCut,
		withSpectrumImgGenerated: withSpectrum,
	}

	msg, err := handleBandpassCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func lowpassExecutioner(imgPath string, args map[string]string) ExecutionResult {
	cutoff, err := parseIntArg(args, "cutoff")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		cutoff:                   cutoff,
		withSpectrumImgGenerated: withSpectrum,
	}

	msg, err := handleLowpassCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}

}

func highpassExecutioner(imgPath string, args map[string]string) ExecutionResult {
	cutoff, err := parseIntArg(args, "cutoff")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		cutoff:                   cutoff,
		withSpectrumImgGenerated: withSpectrum,
	}

	msg, err := handleHighpassCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func bandcutExecutioner(imgPath string, args map[string]string) ExecutionResult {
	lowCut, err := parseIntArg(args, "lowCut")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}
	highCut, err := parseIntArg(args, "highCut")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}
	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		lowCut:                   lowCut,
		highCut:                  highCut,
		withSpectrumImgGenerated: withSpectrum,
	}

	msg, err := handleBandcutCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func phasemodExecutioner(imgPath string, args map[string]string) ExecutionResult {
	k, err := parseIntArg(args, "k")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	l, err := parseIntArg(args, "l")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath: imgPath,
		k:       k,
		l:       l,
	}

	msg, err := handlePhasemodCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func maskpassExecutioner(imgPath string, args map[string]string) ExecutionResult {
	maskName := strings.TrimSpace(args["maskName"])
	if maskName == "" {
		return ExecutionResult{
			Message: "",
			Err:     errors.New("mask name cannot be empty"),
		}
	}

	maskPath := filepath.Join("orthogonal_transforms", "masks", maskName)

	withSpectrum, err := parseBoolArg(args, "withSpectrumImgGenerated")
	if err != nil {
		return ExecutionResult{
			Message: "",
			Err:     err,
		}
	}

	opts := handlingCommandOptions{
		imgPath:                  imgPath,
		maskPath:                 maskPath,
		withSpectrumImgGenerated: withSpectrum,
	}

	msg, err := handleMaskpassCommand(opts)

	return ExecutionResult{
		Message: msg,
		Err:     err,
	}
}

func ExecuteCommand(imgPath, cmdName string, cmdArgs map[string]string) ExecutionResult {
	handler, exists := commandRegistry[cmdName]
	if !exists {
		return ExecutionResult{
			Message: "",
			Err:     fmt.Errorf("command not found: %s", cmdName),
		}
	}
	return handler(imgPath, cmdArgs)
}
