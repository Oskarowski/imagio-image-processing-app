package orthogonal_transforms

import (
	"image-processing/imageio"
	"image-processing/morphological"
	"log"
	"testing"
)

func TestCameraHighpass(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/camera.bmp")
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
	imageio.SaveBmpImage(magnitudeImg, "test_camera_highpass_30_normalized_magnitude_spectrum.bmp")

	filteredSpectrum := HighPassFilter2D(shiftedSpectrum, 30)

	filteredSpectrum[0][0] = dcComponent

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)

	magnitude2 := FFTMagnitudeSpectrum(unshiftedFilteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_camera_highpass_30_normalized_unshifted_after_filtered_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_camera_highpass_30_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_camera_highpass_30_visualize_image.bmp")

}

func TestHighPassFilterWithEdgeDetection2DImage2(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/F5test2.bmp")
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
	imageio.SaveBmpImage(magnitudeImg, "test_F5test2_High_Pass_With_Edge_Detection_normalized_magnitude_spectrum.bmp")

	maskImg, err := imageio.OpenBmpImage("./masks/F5mask2.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	binaryMask := morphological.ConvertIntoBinaryImage(maskImg)

	filteredSpectrum := HighPassFilterWithEdgeDetection2D(shiftedSpectrum, binaryMask)

	magnitudeAfterFilter := FFTMagnitudeSpectrum(filteredSpectrum)
	normalizedAfterFilter := NormalizeMagnitude(magnitudeAfterFilter)
	magnitudeImgAfterFilter := MagnitudeToImage(normalizedAfterFilter)
	imageio.SaveBmpImage(magnitudeImgAfterFilter, "test_F5test2_High_Pass_With_Edge_Detection_after_filtered_magnitude_spectrum.bmp")

	filteredSpectrum[0][0] = dcComponent

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)

	magnitude2 := FFTMagnitudeSpectrum(unshiftedFilteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_F5test2_High_Pass_With_Edge_Detection_normalized_unshifted_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_F5test2_High_Pass_With_Edge_Detection_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_F5test2_High_Pass_With_Edge_Detection_visualize_image.bmp")
}

func TestBoatBandpass(t *testing.T) {
	loadedImg, err := imageio.OpenBmpImage("../imgs/boat.bmp")
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
	imageio.SaveBmpImage(magnitudeImg, "test_boat_bandpass_magnitude_spectrum.bmp")

	filteredSpectrum := BandPassFilter2D(shiftedSpectrum, 10.0, 50.0)

	filteredSpectrum[0][0] = dcComponent

	unshiftedFilteredSpectrum := QuadrantsSwap(filteredSpectrum)

	magnitude2 := FFTMagnitudeSpectrum(unshiftedFilteredSpectrum)
	normalized2 := NormalizeMagnitude(magnitude2)
	magnitudeImg2 := MagnitudeToImage(normalized2)
	imageio.SaveBmpImage(magnitudeImg2, "test_boat_bandpass_unshifted_after_filtered_magnitude_spectrum.bmp")

	reconstructedImageMatrix := FFT2D(unshiftedFilteredSpectrum, true)

	imageio.SaveBmpImage(ConvertComplexToImage(reconstructedImageMatrix), "test_boat_bandpass_converted_image.bmp")
	imageio.SaveBmpImage(ConvertFloatMatrixToImage(VisualizeSpectrum(reconstructedImageMatrix)), "test_boat_bandpass_30_visualize_image.bmp")

}
