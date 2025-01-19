package orthogonal_transforms

import (
	"image"
	"image-processing/morphological"
	"strconv"
)

type SpectrumImage struct {
	Img  image.RGBA
	Name string
}

func HandleBandpassFiltering(img image.Image, filename string, lowCut, highCut int, withSpectrum bool) []SpectrumImage {
	complexMatrix := ConvertImageToComplex(img)
	ftm := FFT2D(complexMatrix, false)

	dcComponent := ftm[0][0]
	shiftedFtm := QuadrantsSwap(ftm)

	generatedImgs := make([]SpectrumImage, 0)

	if withSpectrum {
		magnitude := FFTMagnitudeSpectrum(shiftedFtm)
		normalized := NormalizeMagnitude(magnitude)
		magnitudeImg := MagnitudeToImage(normalized)
		fn := filename + "_magnitude_spectrum.bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*magnitudeImg, fn})
	}

	filtered := BandPassFilter2D(shiftedFtm, float64(lowCut), float64(highCut))

	unshiftedFiltered := QuadrantsSwap(filtered)
	unshiftedFiltered[0][0] = dcComponent

	if withSpectrum {
		mf := FFTMagnitudeSpectrum(filtered)
		nf := NormalizeMagnitude(mf)
		mg := MagnitudeToImage(nf)
		fn := filename + "_bandpass_filtered_magnitude_spectrum_f_" + strconv.Itoa(lowCut) + "_t_" + strconv.Itoa(highCut) + ".bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*mg, fn})
	}

	iftm := FFT2D(unshiftedFiltered, true)
	filteredImg := ConvertComplexToImage(iftm)

	fn := filename + "_bandpass_f_" + strconv.Itoa(lowCut) + "_t_" + strconv.Itoa(highCut) + ".bmp"
	generatedImgs = append(generatedImgs, SpectrumImage{*filteredImg, fn})

	return generatedImgs
}

func HandleLowpassFiltering(img image.Image, filename string, cutoff int, withSpectrum bool) []SpectrumImage {
	complexMatrix := ConvertImageToComplex(img)
	ftm := FFT2D(complexMatrix, false)

	dcComponent := ftm[0][0]
	shiftedFtm := QuadrantsSwap(ftm)

	generatedImgs := make([]SpectrumImage, 0)

	if withSpectrum {
		magnitude := FFTMagnitudeSpectrum(shiftedFtm)
		normalized := NormalizeMagnitude(magnitude)
		magnitudeImg := MagnitudeToImage(normalized)
		fn := filename + "_magnitude_spectrum.bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*magnitudeImg, fn})
	}

	filtered := LowPassFilter2D(shiftedFtm, float64(cutoff))

	unshiftedFiltered := QuadrantsSwap(filtered)
	unshiftedFiltered[0][0] = dcComponent

	if withSpectrum {
		mf := FFTMagnitudeSpectrum(filtered)
		nf := NormalizeMagnitude(mf)
		mg := MagnitudeToImage(nf)
		fn := filename + "_lowpass_filtered_magnitude_spectrum_cutoff_" + strconv.Itoa(cutoff) + ".bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*mg, fn})
	}

	iftm := FFT2D(unshiftedFiltered, true)
	filteredImg := ConvertComplexToImage(iftm)

	fn := filename + "_lowpass_cutoff_" + strconv.Itoa(cutoff) + ".bmp"
	generatedImgs = append(generatedImgs, SpectrumImage{*filteredImg, fn})

	return generatedImgs
}

func HandleHighpassFiltering(img image.Image, filename string, cutoff int, withSpectrum bool) []SpectrumImage {
	complexMatrix := ConvertImageToComplex(img)
	ftm := FFT2D(complexMatrix, false)

	dcComponent := ftm[0][0]
	shiftedFtm := QuadrantsSwap(ftm)

	generatedImgs := make([]SpectrumImage, 0)

	if withSpectrum {
		magnitude := FFTMagnitudeSpectrum(shiftedFtm)
		normalized := NormalizeMagnitude(magnitude)
		magnitudeImg := MagnitudeToImage(normalized)
		fn := filename + "_magnitude_spectrum.bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*magnitudeImg, fn})
	}

	filtered := HighPassFilter2D(shiftedFtm, float64(cutoff))

	unshiftedFiltered := QuadrantsSwap(filtered)
	unshiftedFiltered[0][0] = dcComponent

	if withSpectrum {
		mf := FFTMagnitudeSpectrum(filtered)
		nf := NormalizeMagnitude(mf)
		mg := MagnitudeToImage(nf)
		fn := filename + "_highpass_filtered_magnitude_spectrum_cutoff_" + strconv.Itoa(cutoff) + ".bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*mg, fn})
	}

	iftm := FFT2D(unshiftedFiltered, true)
	filteredImg := ConvertComplexToImage(iftm)

	fn := filename + "_highpass_cutoff_" + strconv.Itoa(cutoff) + ".bmp"
	generatedImgs = append(generatedImgs, SpectrumImage{*filteredImg, fn})

	return generatedImgs
}

func HandleBandcutFiltering(img image.Image, filename string, lowCut, highCut int, withSpectrum bool) []SpectrumImage {
	complexMatrix := ConvertImageToComplex(img)
	ftm := FFT2D(complexMatrix, false)

	dcComponent := ftm[0][0]
	shiftedFtm := QuadrantsSwap(ftm)

	generatedImgs := make([]SpectrumImage, 0)

	if withSpectrum {
		magnitude := FFTMagnitudeSpectrum(shiftedFtm)
		normalized := NormalizeMagnitude(magnitude)
		magnitudeImg := MagnitudeToImage(normalized)
		fn := filename + "_magnitude_spectrum.bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*magnitudeImg, fn})
	}

	filtered := BandCutFilter2D(shiftedFtm, float64(lowCut), float64(highCut))

	unshiftedFiltered := QuadrantsSwap(filtered)
	unshiftedFiltered[0][0] = dcComponent

	if withSpectrum {
		mf := FFTMagnitudeSpectrum(filtered)
		nf := NormalizeMagnitude(mf)
		mg := MagnitudeToImage(nf)
		fn := filename + "_bandcut_filtered_magnitude_spectrum_cutoff_f_" + strconv.Itoa(lowCut) + "_t_" + strconv.Itoa(highCut) + ".bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*mg, fn})
	}

	iftm := FFT2D(unshiftedFiltered, true)
	filteredImg := ConvertComplexToImage(iftm)

	fn := filename + "_bandcut_cutoff_f_" + strconv.Itoa(lowCut) + "_t_" + strconv.Itoa(highCut) + ".bmp"
	generatedImgs = append(generatedImgs, SpectrumImage{*filteredImg, fn})

	return generatedImgs
}

func HandlePhaseModification(img image.Image, filename string, k, l int) []SpectrumImage {
	complexMatrix := ConvertImageToComplex(img)
	ftm := FFT2D(complexMatrix, false)

	shiftedFtm := QuadrantsSwap(ftm)

	generatedImgs := make([]SpectrumImage, 0)

	filtered := PhaseModifyingFilter(shiftedFtm, k, l)

	unshiftedFiltered := QuadrantsSwap(filtered)

	iftm := FFT2D(unshiftedFiltered, true)
	filteredImg := ConvertComplexToImage(iftm)

	fn := filename + "_phase_modified_k_" + strconv.Itoa(k) + "_l_" + strconv.Itoa(l) + ".bmp"
	generatedImgs = append(generatedImgs, SpectrumImage{*filteredImg, fn})

	return generatedImgs
}

func HandleMaskpassFiltering(img image.Image, filename string, maskImg image.Image, withSpectrum bool) []SpectrumImage {
	complexMatrix := ConvertImageToComplex(img)
	ftm := FFT2D(complexMatrix, false)

	dcComponent := ftm[0][0]
	shiftedFtm := QuadrantsSwap(ftm)

	generatedImgs := make([]SpectrumImage, 0)

	if withSpectrum {
		magnitude := FFTMagnitudeSpectrum(shiftedFtm)
		normalized := NormalizeMagnitude(magnitude)
		magnitudeImg := MagnitudeToImage(normalized)
		fn := filename + "_magnitude_spectrum.bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*magnitudeImg, fn})
	}

	binaryMask := morphological.ConvertIntoBinaryImage(maskImg)

	filtered := HighPassFilterWithEdgeDetection2D(shiftedFtm, binaryMask)

	unshiftedFiltered := QuadrantsSwap(filtered)
	unshiftedFiltered[0][0] = dcComponent

	if withSpectrum {
		mf := FFTMagnitudeSpectrum(filtered)
		nf := NormalizeMagnitude(mf)
		mg := MagnitudeToImage(nf)
		fn := filename + "_maskpass_filtered_magnitude_spectrum.bmp"
		generatedImgs = append(generatedImgs, SpectrumImage{*mg, fn})
	}

	iftm := FFT2D(unshiftedFiltered, true)
	filteredImg := ConvertComplexToImage(iftm)

	fn := filename + "_maskpass.bmp"
	generatedImgs = append(generatedImgs, SpectrumImage{*filteredImg, fn})

	return generatedImgs
}
