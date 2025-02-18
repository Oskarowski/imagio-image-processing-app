package orthogonal_transforms

import (
	"imagio/imageio"
	"imagio/morphological"
	"log"
	"testing"
)

func TestBandpassFilter(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/mandril.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	dcComponent := ftSpectrum[0][0]

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_bandpass_magnitude_spectrum.bmp")

	filteredSpectrum := BandPassFilter2D(shiftedSpectrum, 25.0, 90.0)

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)

	unshiftedFilteredSpectrum[0][0] = dcComponent

	magnitude2 := FFTMagnitudeSpectrum(filteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_bandpass_unshifted_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_bandpass_converted_image.bmp")
}

func TestLowpassFilter(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/lenag.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	dcComponent := ftSpectrum[0][0]

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_lowpass_magnitude_spectrum.bmp")

	filteredSpectrum := LowPassFilter2D(shiftedSpectrum, 5.0)

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)

	unshiftedFilteredSpectrum[0][0] = dcComponent

	magnitude2 := FFTMagnitudeSpectrum(filteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_lowpass_unshifted_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_lowpass_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_phase_lowpass_visualize_image.bmp")
}

func TestHighpassFilter(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/pentagon.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	dcComponent := ftSpectrum[0][0]

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_highpass_magnitude_spectrum.bmp")

	filteredSpectrum := HighPassFilter2D(shiftedSpectrum, 45.0)

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)
	unshiftedFilteredSpectrum[0][0] = dcComponent

	magnitude2 := FFTMagnitudeSpectrum(filteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_highpass_unshifted_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_highpass_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_highpass_visualize_image.bmp")
}

func TestBandCutFilter(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/messer.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	dcComponent := ftSpectrum[0][0]

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_bandcut_magnitude_spectrum.bmp")

	filteredSpectrum := BandCutFilter2D(shiftedSpectrum, 35, 100)

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)
	unshiftedFilteredSpectrum[0][0] = dcComponent

	magnitude2 := FFTMagnitudeSpectrum(filteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_bandcut_unshifted_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_bandcut_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_bandcut_visualize_image.bmp")

}

func TestHighPassFilterWithEdgeDetection2DImage2(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/F5test3.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	dcComponent := ftSpectrum[0][0]

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_highpass_edge_detection_magnitude_spectrum.bmp")

	maskImg, err := imageio.OpenBmpImage("./masks/F5mask2.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	binaryMask := morphological.ConvertIntoBinaryImage(maskImg)

	filteredSpectrum := HighPassFilterWithEdgeDetection2D(shiftedSpectrum, binaryMask)

	magnitudeAfterFilter := FFTMagnitudeSpectrum(filteredSpectrum)
	normalizedAfterFilter := NormalizeMagnitude(magnitudeAfterFilter)
	magnitudeImgAfterFilter := MagnitudeToImage(normalizedAfterFilter)
	imageio.SaveBmpImage(magnitudeImgAfterFilter, "test_highpass_edge_detection_filtered_magnitude_spectrum.bmp")

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)
	unshiftedFilteredSpectrum[0][0] = dcComponent

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_highpass_edge_detection_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_highpass_edge_detection_visualize_image.bmp")
}
func TestHighPassFilterWithEdgeDetection2DOnLena(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/lenag.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	dcComponent := ftSpectrum[0][0]

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_lena_edge_detection_magnitude_spectrum.bmp")

	maskImg, err := imageio.OpenBmpImage("./masks/F5mask5.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	binaryMask := morphological.ConvertIntoBinaryImage(maskImg)

	filteredSpectrum := HighPassFilterWithEdgeDetection2D(shiftedSpectrum, binaryMask)

	magnitudeAfterFilter := FFTMagnitudeSpectrum(filteredSpectrum)
	normalizedAfterFilter := NormalizeMagnitude(magnitudeAfterFilter)
	magnitudeImgAfterFilter := MagnitudeToImage(normalizedAfterFilter)
	imageio.SaveBmpImage(magnitudeImgAfterFilter, "test_lena_edge_detection_after_filtered_magnitude_spectrum.bmp")

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)
	unshiftedFilteredSpectrum[0][0] = dcComponent

	magnitude2 := FFTMagnitudeSpectrum(unshiftedFilteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_lena_edge_detection_unshifted_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_lena_edge_detection_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_lena_edge_detection_visualize_image.bmp")
}

func TestPhaseModifyingFilter(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/mandril.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	complexImg := ConvertImageToComplex(loadedImg)

	ftSpectrum := FFT2D(complexImg, false)

	shiftedSpectrum := QuadrantsSwap(ftSpectrum)

	magnitude := FFTMagnitudeSpectrum(shiftedSpectrum)
	normalized := NormalizeMagnitude(magnitude)
	magnitudeImg := MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "test_phase_modifying_filter_magnitude_spectrum.bmp")

	k := 125
	l := 125
	filteredSpectrum := PhaseModifyingFilter(shiftedSpectrum, k, l)

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)

	magnitude2 := FFTMagnitudeSpectrum(filteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_phase_modifying_filter_after_filtered_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_phase_modifying_filter_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_phase_modifying_filter_visualize_image.bmp")

}
