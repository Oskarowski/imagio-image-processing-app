package tui

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

func executeCommand(imgPath, cmdName string, cmdArgs map[string]string) (string, error) {

	log.Default().Printf("Executing command commander here: %s", cmdName)

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

	default:
		return "", errors.New("command not found")
	}

}
