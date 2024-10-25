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

type commandInvocation struct {
	Name        string
	Description string
	Result      string
	Duration    time.Duration
}

func main() {
	if len(os.Args) < 3 {
		cmd.PrintHelp()
		return
	}

	if os.Args[1] == "--help" {
		cmd.PrintHelp()
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

	var commandResults []commandInvocation

	for _, command := range commands {
		cmdResult := commandInvocation{Name: command.Name}
		startTime := time.Now()

		switch command.Name {
		case "brightness":
			brightness, err := strconv.Atoi(command.Args["value"])
			if err != nil {
				log.Fatalf("Brightness value must be int number: %v", err)
			}

			newImg = manipulations.AdjustBrightness(img, brightness)
			outputFileName = fmt.Sprintf("%s_altered_brightness.bmp", originalNameWithoutExt)
			cmdResult.Description = fmt.Sprintf("Brightness adjusted by %d", brightness)

		case "contrast":
			contrast, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Contrast value must be int number: %v", err)
			}

			if contrast < -255 || contrast > 255 {
				log.Fatalf("Contrast value must be in the range of -255 to 255")
			}

			newImg = manipulations.AdjustContrast(img, contrast)
			outputFileName = fmt.Sprintf("%s_altered_contrast.bmp", originalNameWithoutExt)
			cmdResult.Description = fmt.Sprintf("Contrast adjusted by %d", contrast)

		case "negative":
			newImg = manipulations.NegativeImage(img)
			outputFileName = fmt.Sprintf("%s_negative.bmp", originalNameWithoutExt)
			cmdResult.Description = "Negative image created"

		case "hflip":
			newImg = manipulations.HorizontalFlip(img)
			outputFileName = fmt.Sprintf("%s_horizontal_flip.bmp", originalNameWithoutExt)
			cmdResult.Description = "Image horizontally flipped"

		case "vflip":
			newImg = manipulations.VerticalFlip(img)
			outputFileName = fmt.Sprintf("%s_vertical_flip.bmp", originalNameWithoutExt)
			cmdResult.Description = "Image vertically flipped"

		case "dflip":
			newImg = manipulations.DiagonalFlip(img)
			outputFileName = fmt.Sprintf("%s_diagonal_flip.bmp", originalNameWithoutExt)
			cmdResult.Description = "Image diagonally flipped"

		case "shrink":
			factor, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Shrink factor value must be int number: %v", err)
			}

			newImg, err = manipulations.ShrinkImage(img, factor)
			if err != nil {
				log.Fatalf("Error shrinking image: %v", err)
			}

			outputFileName = fmt.Sprintf("%s_shrunk_by_%dx.bmp", originalNameWithoutExt, factor)
			cmdResult.Description = fmt.Sprintf("Image shrunk by a factor of %d", factor)

		case "enlarge":
			factor, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Enlarge factor value must be int number: %v", err)

			}

			newImg, err = manipulations.EnlargeImage(img, factor)
			if err != nil {
				log.Fatalf("Error enlarging image: %v", err)
			}

			outputFileName = fmt.Sprintf("%s_enlarged_by_%dx.bmp", originalNameWithoutExt, factor)
			cmdResult.Description = fmt.Sprintf("Image enlarged by a factor of %d", factor)

		case "adaptive":

			minWindowSize := cmd.AtoiOrDefault(command.Args["min"], 3)
			maxWindowSize := cmd.AtoiOrDefault(command.Args["max"], 7)

			if maxWindowSize < minWindowSize {
				log.Fatal("Max window size must be greater than min window size")
			}

			newImg = noise.AdaptiveMedianFilter(img, minWindowSize, maxWindowSize)

			outputFileName = fmt.Sprintf("%s_adaptive_median_filter.bmp", originalNameWithoutExt)
			cmdResult.Description = "Adaptive median filter applied"

		case "min":
			windowSize, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Window size must be an int: %v", err)
			}
			newImg = noise.MinFilter(img, windowSize)
			outputFileName = fmt.Sprintf("%s_min_filter.bmp", originalNameWithoutExt)
			cmdResult.Description = fmt.Sprintf("Min filter applied with window size %d", windowSize)

		case "max":
			windowSize, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Window size must be an int: %v", err)
			}
			newImg = noise.MaxFilter(img, windowSize)
			outputFileName = fmt.Sprintf("%s_max_filter.bmp", originalNameWithoutExt)
			cmdResult.Description = fmt.Sprintf("Max filter applied with window size %d", windowSize)

		case "mse":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for MSE.")
			}

			var mse float64

			if newImg != nil {
				mse = analysis.MeanSquareError(newImg, comparisonImage)
			} else {
				mse = analysis.MeanSquareError(img, comparisonImage)
			}

			cmdResult.Description = "Mean Square Error calculated"
			cmdResult.Result = fmt.Sprintf("MSE: %f", mse)

		case "pmse":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for PMSE.")
			}

			var pmse float64

			if newImg != nil {
				pmse = analysis.PeakMeanSquareError(newImg, comparisonImage)
			} else {
				pmse = analysis.PeakMeanSquareError(img, comparisonImage)
			}

			cmdResult.Description = "Peak Mean Square Error calculated"
			cmdResult.Result = fmt.Sprintf("PMSE: %f", pmse)

		case "snr":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for SNR.")
			}

			var snr float64

			if newImg != nil {
				snr = analysis.SignalToNoiseRatio(newImg, comparisonImage)
			} else {
				snr = analysis.SignalToNoiseRatio(img, comparisonImage)
			}

			cmdResult.Description = "Signal to Noise Ratio calculated"
			cmdResult.Result = fmt.Sprintf("SNR: %f", snr)

		case "psnr":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for PSNR.")
			}

			var psnr float64

			if newImg != nil {
				psnr = analysis.PeakSignalToNoiseRatio(newImg, comparisonImage)
			} else {
				psnr = analysis.PeakSignalToNoiseRatio(img, comparisonImage)
			}

			cmdResult.Description = "Peak Signal to Noise Ratio calculated"
			cmdResult.Result = fmt.Sprintf("PSNR: %f", psnr)

		case "md":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for MD.")
			}

			var md int

			if newImg != nil {
				md = analysis.MaxDifference(newImg, comparisonImage)
			} else {
				md = analysis.MaxDifference(img, comparisonImage)
			}

			cmdResult.Description = "Max Difference calculated"
			cmdResult.Result = fmt.Sprintf("Max Difference: %d", md)

		default:
			fmt.Println("Unknown commend")
			return
		}

		cmdResult.Duration = time.Since(startTime)
		commandResults = append(commandResults, cmdResult)

		durationSum += cmdResult.Duration
	}

	if newImg != nil {
		err = imageio.SaveBmpImage(newImg, outputFileName)
		if err != nil {
			log.Fatalf("Error saving file: %v", err)
		} else {
			fmt.Printf("Image saved successfully as: %s\n", outputFileName)
		}
	}

	fmt.Println("Execution Report:")
	for _, result := range commandResults {
		fmt.Printf("Command: %s\n", result.Name)
		fmt.Printf("Description: %s\n", result.Description)
		if result.Result != "" {
			fmt.Printf("Result: %s\n", result.Result)
		}
		fmt.Printf("Duration: %v\n\n", result.Duration)
	}

	fmt.Printf("Total operation time: %v\n", durationSum)
}
