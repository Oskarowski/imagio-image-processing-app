package tui

import (
	"errors"
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
			return "", err
		}

		highCutInt, err := strconv.Atoi(highCut)
		if err != nil {
			return "", err
		}

		withSpectrum, err := strconv.ParseBool(withSpectrumImgGenerated)
		if err != nil {
			return "", err
		}

		return handleBandpassCommand(imgPath, lowCutInt, highCutInt, withSpectrum)

	default:
		return "", errors.New("command not found")
	}

}
