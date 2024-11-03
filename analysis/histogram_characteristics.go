package analysis

func CalculateMean(histogram []int) float64 {
	var sum float64
	var totalPixels int

	for brightness, count := range histogram {
		sum += float64(brightness) * float64(count)
		totalPixels += count
	}

	return sum / float64(totalPixels)
}

func CalculateVariance(histogram []int) float64 {
	var varianceSum float64
	var totalPixels int

	mean := CalculateMean(histogram)

	for brightness, count := range histogram {
		diff := float64(brightness) - mean
		varianceSum += diff * diff * float64(count)
		totalPixels += count
	}

	return varianceSum / float64(totalPixels)
}
