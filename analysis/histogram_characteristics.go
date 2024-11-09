package analysis

import (
	"fmt"
	"math"
)

func calculateMean(histogram [256]int) float64 {
	var sum float64
	var totalPixels int

	for intensity, count := range histogram {
		sum += float64(intensity) * float64(count)
		totalPixels += count
	}

	return sum / float64(totalPixels)
}

func calculateVariance(histogram [256]int) float64 {
	var varianceSum float64
	var totalPixels int

	mean := calculateMean(histogram)

	for intensity, count := range histogram {
		diff := float64(intensity) - mean
		varianceSum += diff * diff * float64(count)
		totalPixels += count
	}

	return varianceSum / float64(totalPixels)
}

func calculateStandardDeviation(histogram [256]int) float64 {
	return math.Sqrt(calculateVariance(histogram))
}

func calculateVariationCoefficientOne(histogram [256]int) float64 {
	meanIntensity := calculateMean(histogram)
	standardDeviation := calculateStandardDeviation(histogram)

	return standardDeviation / meanIntensity
}

func calculateAsymmetryCoefficient(histogram [256]int) float64 {
	meanIntensity := calculateMean(histogram)
	standardDeviation := calculateStandardDeviation(histogram)

	if standardDeviation == 0 {
		return 0
	}

	var sumCubedDiffs float64
	var totalPixels int

	for intensity, count := range histogram {
		diff := float64(intensity) - meanIntensity
		sumCubedDiffs += diff * diff * diff * float64(count)
		totalPixels += count
	}

	asymmetryCoefficient := (1 / (standardDeviation * standardDeviation * standardDeviation)) * (sumCubedDiffs / float64(totalPixels))
	return asymmetryCoefficient
}

func calculateFlatteningCoefficient(histogram [256]int) float64 {
	meanIntensity := calculateMean(histogram)
	standardDeviation := calculateStandardDeviation(histogram)

	if standardDeviation == 0 {
		return 0
	}

	var sumFourthPowerDiffs float64
	var totalPixels int

	for intensity, count := range histogram {
		diff := float64(intensity) - meanIntensity
		sumFourthPowerDiffs += diff * diff * diff * diff * float64(count)
		totalPixels += count
	}

	N := float64(totalPixels)
	fourthMoment := sumFourthPowerDiffs / N

	flatteningCoefficient := (fourthMoment / (standardDeviation * standardDeviation * standardDeviation * standardDeviation)) - 3

	return flatteningCoefficient

}

func calculateVariationCoefficientTwo(histogram [256]int) float64 {
	var sumOfSquares float64
	var totalPixels int

	for _, count := range histogram {
		sumOfSquares += float64(count * count)
		totalPixels += count
	}

	if totalPixels == 0 {
		return 0
	}

	variationCoefficientII := (1.0 / float64(totalPixels)) * (1.0 / float64(totalPixels)) * sumOfSquares

	return variationCoefficientII
}

func calculateInformationSourceEntropy(histogram [256]int) float64 {
	var entropy float64
	var totalPixels int

	for _, count := range histogram {
		totalPixels += count
	}

	if totalPixels == 0 {
		return 0
	}

	N := float64(totalPixels)

	for _, count := range histogram {
		if count == 0 {
			continue
		}

		entropy += float64(count) * math.Log2(float64(count)/N)
	}

	return -entropy / N
}

func CalculateHistogramCharacteristic(metricMethod string, providedHistogram [256]int, filenameWithoutExt string) (string, string) {
	var result string
	var description string

	switch metricMethod {
	case "cmean":
		mean := calculateMean(providedHistogram)
		description = fmt.Sprintf("Calculated Mean intensity for %s", filenameWithoutExt)
		result = fmt.Sprintf("Mean: %f", mean)

	case "cvariance":
		variance := calculateVariance(providedHistogram)
		description = fmt.Sprintf("Calculated Variance intensity for %s", filenameWithoutExt)
		result = fmt.Sprintf("Variance: %f", variance)

	case "cstdev":
		stdev := calculateStandardDeviation(providedHistogram)
		description = fmt.Sprintf("Calculated Standard Deviation for %s", filenameWithoutExt)
		result = fmt.Sprintf("Standard Deviation: %f", stdev)

	case "cvarcoi":
		varCoefI := calculateVariationCoefficientOne(providedHistogram)
		description = fmt.Sprintf("Calculated Variation Coefficient I for %s", filenameWithoutExt)
		result = fmt.Sprintf("Variation Coefficient I: %f", varCoefI)

	case "casyco":
		asymCoef := calculateAsymmetryCoefficient(providedHistogram)
		description = fmt.Sprintf("Calculated Asymmetry Coefficient for %s", filenameWithoutExt)
		result = fmt.Sprintf("Asymmetry Coefficient: %f", asymCoef)

	case "cflatco":
		flatCoef := calculateFlatteningCoefficient(providedHistogram)
		description = fmt.Sprintf("Calculated Flattening Coefficient for %s", filenameWithoutExt)
		result = fmt.Sprintf("Flattening Coefficient: %f", flatCoef)

	case "cvarcoii":
		varCoefII := calculateVariationCoefficientTwo(providedHistogram)
		description = fmt.Sprintf("Calculated Variation Coefficient II for %s", filenameWithoutExt)
		result = fmt.Sprintf("Variation Coefficient II: %f", varCoefII)

	case "centropy":
		entropy := calculateInformationSourceEntropy(providedHistogram)
		description = fmt.Sprintf("Calculated Information Source Entropy for %s", filenameWithoutExt)
		result = fmt.Sprintf("Information Source Entropy: %f", entropy)

	default:
		description = "Unknown metric"
		result = "N/A"
	}

	return result, description
}
