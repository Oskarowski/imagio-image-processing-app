package main

import (
	"fmt"
	"image"
	"image-processing/analysis"
	"image-processing/cmd"
	"image-processing/imageio"
	"image-processing/manipulations"
	"image-processing/noise"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <command> [-argument=value [...]] <bmp_image_path> [<second_image_path>]")
		return
	}

	imagePath := os.Args[len(os.Args)-1]

	var comparisonImagePath string
	if len(os.Args) > 3 && cmd.IsImagePath(os.Args[len(os.Args)-2]) {
		comparisonImagePath = os.Args[len(os.Args)-2]
	}

	img, err := imageio.OpenBmpImage(imagePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	var comparisonImage image.Image
	if comparisonImagePath != "" {
		comparisonImage, err = imageio.OpenBmpImage(comparisonImagePath)
		if err != nil {
			log.Fatalf("Error opening comparison image: %v", err)
		}
	}

	commands := cmd.ParseCommands(os.Args[1 : len(os.Args)-1])

	originalName := filepath.Base(imagePath)
	originalNameWithoutExt := originalName[:len(originalName)-len(filepath.Ext(originalName))]

	var outputFileName string
	var newImg *image.RGBA
	var durationSum time.Duration

	for _, cmd := range commands {
		startTime := time.Now()

		switch cmd.Name {
		case "brightness":
			brightness, err := strconv.Atoi(cmd.Args["value"])
			if err != nil {
				log.Fatalf("Brightness value must be int number: %v", err)
			}

			newImg = manipulations.AdjustBrightness(img, brightness)
			outputFileName = fmt.Sprintf("%s_altered_brightness.bmp", originalNameWithoutExt)

		case "contrast":
			contrast, err := strconv.Atoi(cmd.Args["value"])

			if err != nil {
				log.Fatalf("Contrast value must be int number: %v", err)
			}

			if contrast < -255 || contrast > 255 {
				log.Fatalf("Contrast value must be in the range of -255 to 255")
			}

			newImg = manipulations.AdjustContrast(img, contrast)
			outputFileName = fmt.Sprintf("%s_altered_contrast.bmp", originalNameWithoutExt)

		case "negative":
			newImg = manipulations.NegativeImage(img)
			outputFileName = fmt.Sprintf("%s_negative.bmp", originalNameWithoutExt)

		case "hflip":
			newImg = manipulations.HorizontalFlip(img)
			outputFileName = fmt.Sprintf("%s_horizontal_flip.bmp", originalNameWithoutExt)

		case "vflip":
			newImg = manipulations.VerticalFlip(img)
			outputFileName = fmt.Sprintf("%s_vertical_flip.bmp", originalNameWithoutExt)

		case "dflip":
			newImg = manipulations.DiagonalFlip(img)
			outputFileName = fmt.Sprintf("%s_diagonal_flip.bmp", originalNameWithoutExt)

		case "shrink":
			factor, err := strconv.Atoi(cmd.Args["value"])

			if err != nil {
				log.Fatalf("Shrink factor value must be int number: %v", err)
			}

			newImg, err = manipulations.ShrinkImage(img, factor)
			if err != nil {
				log.Fatalf("Error shrinking image: %v", err)
			}

			outputFileName = fmt.Sprintf("%s_shrunk_by_%dx.bmp", originalNameWithoutExt, factor)

		case "enlarge":
			factor, err := strconv.Atoi(cmd.Args["value"])

			if err != nil {
				log.Fatalf("Enlarge factor value must be int number: %v", err)
			}

			newImg, err = manipulations.EnlargeImage(img, factor)
			if err != nil {
				log.Fatalf("Error enlarging image: %v", err)
			}

			outputFileName = fmt.Sprintf("%s_enlarged_by_%dx.bmp", originalNameWithoutExt, factor)

		case "adaptive":
			newImg = noise.AdaptiveMedianFilter(img, 30)

			outputFileName = fmt.Sprintf("%s_adaptive_median_filter.bmp", originalNameWithoutExt)

		case "min":
			windowSize, err := strconv.Atoi(cmd.Args["value"])

			if err != nil {
				log.Fatalf("Window size must be an int: %v", err)
			}
			newImg = noise.MinFilter(img, windowSize)
			outputFileName = fmt.Sprintf("%s_min_filter.bmp", originalNameWithoutExt)

		case "max":
			windowSize, err := strconv.Atoi(cmd.Args["value"])

			if err != nil {
				log.Fatalf("Window size must be an int: %v", err)
			}
			newImg = noise.MaxFilter(img, windowSize)
			outputFileName = fmt.Sprintf("%s_max_filter.bmp", originalNameWithoutExt)

		case "mse":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for MSE.")
			}

			mse := analysis.MeanSquareError(img, comparisonImage)
			fmt.Printf("MSE: %v\n", mse)

		case "pmse":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for PMSE.")
			}

			pmse := analysis.PeakMeanSquareError(img, comparisonImage)
			fmt.Printf("PMSE: %v\n", pmse)

		case "snr":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for SNR.")
			}

			snr := analysis.SignalToNoiseRatio(img, comparisonImage)
			fmt.Printf("SNR: %v\n", snr)

		case "psnr":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for PSNR.")
			}

			psnr := analysis.PeakSignalToNoiseRatio(img, comparisonImage)
			fmt.Printf("PSNR: %v\n", psnr)

		case "md":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for MD.")
			}

			md := analysis.PeakSignalToNoiseRatio(img, comparisonImage)
			fmt.Printf("Maximum Difference: %v\n", md)

		default:
			fmt.Println("Unknown commend")
			return
		}

		duration := time.Since(startTime)
		durationSum += duration
		fmt.Printf("Operation '%s' took: %v\n", cmd.Name, duration)
	}

	if newImg != nil {
		err = imageio.SaveBmpImage(newImg, outputFileName)
		if err != nil {
			log.Fatalf("Error saving file: %v", err)
		} else {
			fmt.Printf("Image saved successfully as: %s\n", outputFileName)
		}
	}

	fmt.Printf("Total operation time: %v\n", durationSum)
}
