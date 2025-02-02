package tui

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func executeCommand(imgPath, cmdName string, cmdArgs map[string]string) (string, error) {

	switch cmdName {

	case "bandpass":
		lowCut, highCut, withSpectrumImgGenerated := cmdArgs["lowCut"], cmdArgs["highCut"], cmdArgs["withSpectrumImgGenerated"]

		lowCutInt, err := strconv.Atoi(lowCut)
		if err != nil {
			return "", fmt.Errorf("invalid lowCut value: %w", err)
		}

		highCutInt, err := strconv.Atoi(highCut)
		if err != nil {
			return "", fmt.Errorf("invalid highCut value: %w", err)
		}

		withSpectrum, err := strconv.ParseBool(withSpectrumImgGenerated)
		if err != nil {
			return "", fmt.Errorf("invalid withSpectrumImgGenerated value: %w", err)
		}

		result, err := handleBandpassCommand(imgPath, lowCutInt, highCutInt, withSpectrum)
		if err != nil {
			return "", fmt.Errorf("bandpass filtering failed: %w", err)
		}

		return result, nil

	case "lowpass":
		cutoff, withSpectrumImgGenerated := cmdArgs["cutoff"], cmdArgs["withSpectrumImgGenerated"]

		cutoffInt, err := strconv.Atoi(cutoff)
		if err != nil {
			return "", fmt.Errorf("invalid cutoff value: %w", err)
		}

		withSpectrum, err := strconv.ParseBool(withSpectrumImgGenerated)
		if err != nil {
			return "", fmt.Errorf("invalid withSpectrumImgGenerated value: %w", err)
		}

		result, err := handleLowpassCommand(imgPath, cutoffInt, withSpectrum)
		if err != nil {
			return "", fmt.Errorf("lowpass filtering failed: %w", err)
		}

		return result, nil

	case "highpass":
		cutoff, withSpectrumImgGenerated := cmdArgs["cutoff"], cmdArgs["withSpectrumImgGenerated"]

		cutoffInt, err := strconv.Atoi(cutoff)
		if err != nil {
			return "", fmt.Errorf("invalid cutoff value: %w", err)
		}

		withSpectrum, err := strconv.ParseBool(withSpectrumImgGenerated)
		if err != nil {
			return "", fmt.Errorf("invalid withSpectrumImgGenerated value: %w", err)
		}

		result, err := handleHighpassCommand(imgPath, cutoffInt, withSpectrum)
		if err != nil {
			return "", fmt.Errorf("highpass filtering failed: %w", err)
		}

		return result, nil

	case "bandcut":
		lowCut, highCut, withSpectrumImgGenerated := cmdArgs["lowCut"], cmdArgs["highCut"], cmdArgs["withSpectrumImgGenerated"]

		lowCutInt, err := strconv.Atoi(lowCut)
		if err != nil {
			return "", fmt.Errorf("invalid lowCut value: %w", err)
		}

		highCutInt, err := strconv.Atoi(highCut)
		if err != nil {
			return "", fmt.Errorf("invalid highCut value: %w", err)
		}

		withSpectrum, err := strconv.ParseBool(withSpectrumImgGenerated)
		if err != nil {
			return "", fmt.Errorf("invalid withSpectrumImgGenerated value: %w", err)
		}

		result, err := handleBandcutCommand(imgPath, lowCutInt, highCutInt, withSpectrum)
		if err != nil {
			return "", fmt.Errorf("bandcut filtering failed: %w", err)
		}

		return result, nil

	case "phasemod":
		k, l := cmdArgs["k"], cmdArgs["l"]

		kInt, err := strconv.Atoi(k)
		if err != nil {
			return "", fmt.Errorf("invalid k value: %w", err)
		}

		lInt, err := strconv.Atoi(l)
		if err != nil {
			return "", fmt.Errorf("invalid l value: %w", err)
		}

		result, err := handlePhasemodCommand(imgPath, kInt, lInt)
		if err != nil {
			return "", fmt.Errorf("phase modification failed: %w", err)
		}

		return result, nil

	case "maskpass":
		maskName, withSpectrumImgGenerated := cmdArgs["maskName"], cmdArgs["withSpectrumImgGenerated"]

		maskName = strings.Trim(maskName, " ")
		if maskName == "" {
			return "", errors.New("maskName cannot be empty")
		}
		maskPath := filepath.Join("orthogonal_transforms", "masks", maskName)

		withSpectrum, err := strconv.ParseBool(withSpectrumImgGenerated)
		if err != nil {
			return "", fmt.Errorf("invalid withSpectrumImgGenerated value: %w", err)
		}

		result, err := handleMaskpassCommand(imgPath, maskPath, withSpectrum)
		if err != nil {
			return "", fmt.Errorf("maskpass filtering failed: %w", err)
		}

		return result, nil

	default:
		return "", errors.New("command not found")
	}

}
