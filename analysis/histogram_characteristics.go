package analysis

import "math"

func CalculateMean(histogram [256]int) float64 {
	var sum float64
	var totalPixels int

	for intensity, count := range histogram {
		sum += float64(intensity) * float64(count)
		totalPixels += count
	}

	return sum / float64(totalPixels)
}

func CalculateVariance(histogram [256]int) float64 {
	var varianceSum float64
	var totalPixels int

	mean := CalculateMean(histogram)

	for intensity, count := range histogram {
		diff := float64(intensity) - mean
		varianceSum += diff * diff * float64(count)
		totalPixels += count
	}

	return varianceSum / float64(totalPixels)
}

func CalculateStandardDeviation(histogram [256]int) float64 {
	return math.Sqrt(CalculateVariance(histogram))
}

func CalculateVariationCoefficientOne(histogram [256]int) float64 {
	meanIntensity := CalculateMean(histogram)
	standardDeviation := CalculateStandardDeviation(histogram)

	return standardDeviation / meanIntensity
}

func CalculateAsymmetryCoefficient(histogram [256]int) float64 {
	meanIntensity := CalculateMean(histogram)
	standardDeviation := CalculateStandardDeviation(histogram)

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
