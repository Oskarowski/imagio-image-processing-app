package main

import (
	"image-processing/imageio"
	"image-processing/orthogonal_transforms"
	"log"
)

func main() {
	loadedImg, err := imageio.OpenBmpImage("imgs/boat.bmp")
	// loadedImg, err := imageio.OpenBmpImage("imgs/lenag.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	complexImg := orthogonal_transforms.ConvertImageToComplex(loadedImg)

	// imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(complexImg), "loaded_img_to_complex_to_img.bmp")
	// imageio.SaveBmpImage(orthogonal_transforms.VisualizeSpectrumInImage(complexImg, 512, 512), "visualize_loaded_img.bmp")

	fftData := orthogonal_transforms.FFT2D(complexImg, false)
	dcComponent := fftData[0][0]

	// centeredFrequency := orthogonal_transforms.SwapQuadrants(fftData)

	shiftedFreqDomain := orthogonal_transforms.QuadrantsSwap(fftData)

	magnitude := orthogonal_transforms.FFTMagnitudeSpectrum(shiftedFreqDomain)
	normalized := orthogonal_transforms.NormalizeMagnitude(magnitude)
	magnitudeImg := orthogonal_transforms.MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "magnitude_spectrum.bmp")

	// imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(fftData), "fft_data_of_complex_img.bmp")
	// imageio.SaveBmpImage(orthogonal_transforms.VisualizeSpectrumInImage(fftData, 512, 512), "visualize_fft_data_of_complex_img.bmp")

	// tifftData := orthogonal_transforms.IFFT2D(fftData)
	// imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(tifftData), "inverse_fft_on_fft_data.bmp")

	lowCutoff := 10.0
	highCutoff := 50.0

	bandpassFiltered := orthogonal_transforms.BandPassFilter2D(shiftedFreqDomain, lowCutoff, highCutoff)
	bandpassFiltered[0][0] = dcComponent

	imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(bandpassFiltered), "band_pass_matrix_to_img.bmp")
	imageio.SaveBmpImage(orthogonal_transforms.VisualizeSpectrumInImage(bandpassFiltered, 512, 512), "visualize_band_pass_matrix.bmp")
	imageio.SaveBmpImage(orthogonal_transforms.ConvertFloatMatrixToImage(orthogonal_transforms.VisualizeSpectrum(bandpassFiltered)), "visualize_spectrum_band_pass_matrix.bmp")

	// centeredFrequency2 := orthogonal_transforms.SwapQuadrants(bandpassMatrix)

	unshiftedFreqDomain := orthogonal_transforms.QuadrantsSwap(bandpassFiltered)

	reconstructedImageMatrix := orthogonal_transforms.FFT2D(unshiftedFreqDomain, true)

	imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(reconstructedImageMatrix), "converted_band_pass_filter_img.bmp")
	imageio.SaveBmpImage(orthogonal_transforms.VisualizeSpectrumInImage(reconstructedImageMatrix, 512, 512), "visualize_band_pass_filter_img.bmp")
	imageio.SaveBmpImage(orthogonal_transforms.ConvertFloatMatrixToImage(orthogonal_transforms.VisualizeSpectrum(reconstructedImageMatrix)), "visualize_spectrum_band_pass_filter_img.bmp")

	return

	// if len(os.Args) > 1 && os.Args[1] == "--help" {
	// 	cmd.PrintHelp()
	// 	return
	// }

	// logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatalf("Failed to open log file: %v", err)
	// }
	// defer logFile.Close()
	// log.SetOutput(logFile)

	// if len(os.Args) > 1 {
	// 	cmd.RunAsCliApp()
	// } else {
	// 	gui.RunAsTUIApp()
	// }
}
