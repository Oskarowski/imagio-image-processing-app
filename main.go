package main

import (
	"fmt"
	"image-processing/imageio"
	"image-processing/morphological"
	"image-processing/orthogonal_transforms"
	"log"
)

func main() {
	loadedImg, err := imageio.OpenBmpImage("imgs/F5test1.bmp")
	// loadedImg, err := imageio.OpenBmpImage("imgs/pentagon.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	complexImg := orthogonal_transforms.ConvertImageToComplex(loadedImg)

	imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(complexImg), "loaded_img_to_complex_to_img.bmp")

	ftSpectrum := orthogonal_transforms.FFT2D(complexImg, false)
	fmt.Println("Computed FT data")

	dcComponent := ftSpectrum[0][0]

	shiftedFtSpectrum := orthogonal_transforms.QuadrantsSwap(ftSpectrum)

	magnitude := orthogonal_transforms.FFTMagnitudeSpectrum(shiftedFtSpectrum)
	normalized := orthogonal_transforms.NormalizeMagnitude(magnitude)
	magnitudeImg := orthogonal_transforms.MagnitudeToImage(normalized)
	imageio.SaveBmpImage(magnitudeImg, "magnitude_spectrum.bmp")
	fmt.Println("Computed Magnitude Visualization")

	maskImg, err := imageio.OpenBmpImage("orthogonal_transforms/masks/F5mask1.bmp")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	binaryMask := morphological.ConvertIntoBinaryImage(maskImg)

	filteredSpectrum := orthogonal_transforms.HighPassFilterWithEdgeDetection2D(shiftedFtSpectrum, binaryMask)

	filteredSpectrum[0][0] = dcComponent
	fmt.Println("Computed Filter")

	// imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(bandpassFiltered), "band_pass_matrix_to_img.bmp")
	// imageio.SaveBmpImage(orthogonal_transforms.VisualizeSpectrumInImage(bandpassFiltered, 512, 512), "visualize_band_pass_matrix.bmp")
	// imageio.SaveBmpImage(orthogonal_transforms.ConvertFloatMatrixToImage(orthogonal_transforms.VisualizeSpectrum(bandpassFiltered)), "visualize_spectrum_band_pass_matrix.bmp")

	// centeredFrequency2 := orthogonal_transforms.SwapQuadrants(bandpassMatrix)

	unshiftedFreqDomain := orthogonal_transforms.QuadrantsSwap(filteredSpectrum)

	reconstructedImageMatrix := orthogonal_transforms.FFT2D(unshiftedFreqDomain, true)
	fmt.Println("Computed Inverse FT")

	imageio.SaveBmpImage(orthogonal_transforms.ConvertComplexToImage(reconstructedImageMatrix), "converted_filter_img.bmp")
	// imageio.SaveBmpImage(orthogonal_transforms.VisualizeSpectrumInImage(reconstructedImageMatrix, 512, 512), "visualize_filter_img.bmp")
	imageio.SaveBmpImage(orthogonal_transforms.ConvertFloatMatrixToImage(orthogonal_transforms.VisualizeSpectrum(reconstructedImageMatrix)), "visualize_filter_spectrum_img.bmp")

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
