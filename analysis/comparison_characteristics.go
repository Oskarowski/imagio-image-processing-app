package analysis

import (
	"fmt"
	"image"
	"strings"
)

type ImageComparisonEntry struct {
	MetricMethod string
	Description  string
	Result       string
	Img1Name     string
	Img2Name     string
}

func CalculateComparisonCharacteristic(metricMethod string, img1, img2 image.Image) ImageComparisonEntry {
	var result string
	var description string

	lowerMetricMethod := strings.ToLower(strings.Trim(metricMethod, " "))

	switch lowerMetricMethod {
	case "mse":
		mse := MeanSquareError(img1, img2)
		description = "Mean Square Error calculated"
		result = fmt.Sprintf("MSE: %f", mse)

	case "pmse":
		pmse := PeakMeanSquareError(img1, img2)
		description = "Peak Mean Square Error calculated"
		result = fmt.Sprintf("PMSE: %f", pmse)

	case "snr":
		snr := SignalToNoiseRatio(img1, img2)
		description = "Signal to Noise Ratio calculated"
		result = fmt.Sprintf("SNR: %f", snr)

	case "psnr":
		psnr := PeakSignalToNoiseRatio(img1, img2)
		description = "Peak Signal to Noise Ratio calculated"
		result = fmt.Sprintf("PSNR: %f", psnr)

	case "md":
		md := MaxDifference(img1, img2)
		description = "Max Difference calculated"
		result = fmt.Sprintf("Max Difference: %d", md)

	default:
		description = "Unknown metric"
		result = "N/A"
	}

	return ImageComparisonEntry{
		MetricMethod: strings.ToUpper(metricMethod),
		Description:  description,
		Result:       result,
	}
}
