package cmd

import (
	"fmt"
	"image"
	"image-processing/analysis"
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

type ImageQueueItem struct {
	Image       *image.RGBA
	Filename    string
	Denoised    bool
	IsHistogram bool
}

func RunAsCliApp() {

	imagePath := os.Args[len(os.Args)-1]

	var comparisonImagePath string
	if len(os.Args) > 3 && IsImagePath(os.Args[len(os.Args)-2]) {
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

	commands := ParseCommands(os.Args[1 : len(os.Args)-1])

	loadedMasks, err := manipulations.LoadMasksFromJSON("masks.json")
	if err != nil {
		log.Fatalf("Error loading masks: %v", err)
	}

	originalName := filepath.Base(imagePath)
	originalNameWithoutExt := originalName[:len(originalName)-len(filepath.Ext(originalName))]

	var durationSum time.Duration

	var commandResults []commandInvocation
	var imageQueue []ImageQueueItem

	for _, command := range commands {
		cmdResult := commandInvocation{Name: command.Name}
		startTime := time.Now()

		switch command.Name {
		case "brightness":
			brightness, err := strconv.Atoi(command.Args["value"])
			if err != nil {
				log.Fatalf("Brightness value must be int number: %v", err)
			}

			newImg := manipulations.AdjustBrightness(img, brightness)
			outputFileName := fmt.Sprintf("%s_altered_brightness.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = fmt.Sprintf("Brightness adjusted by %d", brightness)

		case "contrast":
			contrast, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Contrast value must be int number: %v", err)
			}

			if contrast < -255 || contrast > 255 {
				log.Fatalf("Contrast value must be in the range of -255 to 255")
			}

			newImg := manipulations.AdjustContrast(img, contrast)
			outputFileName := fmt.Sprintf("%s_altered_contrast.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = fmt.Sprintf("Contrast adjusted by %d", contrast)

		case "negative":

			outputFileName := fmt.Sprintf("%s_negative.bmp", originalNameWithoutExt)
			newImg := manipulations.NegativeImage(img)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = "Negative image created"

		case "hflip":

			newImg := manipulations.HorizontalFlip(img)
			outputFileName := fmt.Sprintf("%s_horizontal_flip.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = "Image horizontally flipped"

		case "vflip":

			newImg := manipulations.VerticalFlip(img)
			outputFileName := fmt.Sprintf("%s_vertical_flip.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = "Image vertically flipped"

		case "dflip":

			newImg := manipulations.DiagonalFlip(img)
			outputFileName := fmt.Sprintf("%s_diagonal_flip.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = "Image diagonally flipped"

		case "shrink":
			factor, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Shrink factor value must be int number: %v", err)
			}

			newImg, err := manipulations.ShrinkImage(img, factor)
			if err != nil {
				log.Fatalf("Error shrinking image: %v", err)
			}

			outputFileName := fmt.Sprintf("%s_shrunk_by_%dx.bmp", originalNameWithoutExt, factor)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = fmt.Sprintf("Image shrunk by a factor of %d", factor)

		case "enlarge":
			factor, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Enlarge factor value must be int number: %v", err)

			}

			newImg, err := manipulations.EnlargeImage(img, factor)
			if err != nil {
				log.Fatalf("Error enlarging image: %v", err)
			}

			outputFileName := fmt.Sprintf("%s_enlarged_by_%dx.bmp", originalNameWithoutExt, factor)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = fmt.Sprintf("Image enlarged by a factor of %d", factor)

		case "adaptive":

			minWindowSize := GetOrDefault(command.Args["min"], 3)
			maxWindowSize := GetOrDefault(command.Args["max"], 7)

			if maxWindowSize < minWindowSize {
				log.Fatal("Max window size must be greater than min window size")
			}

			newImg := noise.AdaptiveMedianFilter(img, minWindowSize, maxWindowSize)
			outputFileName := fmt.Sprintf("%s_adaptive_median_filter.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName, Denoised: true})

			cmdResult.Description = "Adaptive median filter applied"

		case "adaptive-parallel":

			minWindowSize := GetOrDefault(command.Args["min"], 3)
			maxWindowSize := GetOrDefault(command.Args["max"], 7)

			if maxWindowSize < minWindowSize {
				log.Fatal("Max window size must be greater than min window size")
			}

			newImg := noise.AdaptiveMedianFilterParallel(img, minWindowSize, maxWindowSize)
			outputFileName := fmt.Sprintf("%s_adaptive_parallel_median_filter.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName, Denoised: true})

			cmdResult.Description = fmt.Sprintf("Adaptive median filter applied with min window size %d and max window size %d", minWindowSize, maxWindowSize)

		case "min":
			windowSize, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Window size must be an int: %v", err)
			}
			newImg := noise.MinFilter(img, windowSize)
			outputFileName := fmt.Sprintf("%s_min_filter.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName, Denoised: true})

			cmdResult.Description = fmt.Sprintf("Min filter applied with window size %d", windowSize)

		case "max":
			windowSize, err := strconv.Atoi(command.Args["value"])

			if err != nil {
				log.Fatalf("Window size must be an int: %v", err)
			}
			newImg := noise.MaxFilter(img, windowSize)
			outputFileName := fmt.Sprintf("%s_max_filter.bmp", originalNameWithoutExt)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName, Denoised: true})

			cmdResult.Description = fmt.Sprintf("Max filter applied with window size %d", windowSize)

		case "mse":
			if comparisonImage == nil {
				log.Fatalf("Comparison image is required for MSE.")
			}

			var mse float64
			lastDenoisedImage := getLastDenoisedImage(imageQueue)

			if lastDenoisedImage != nil {
				mse = analysis.MeanSquareError(lastDenoisedImage, comparisonImage)
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
			lastDenoisedImage := getLastDenoisedImage(imageQueue)

			if lastDenoisedImage != nil {
				pmse = analysis.PeakMeanSquareError(lastDenoisedImage, comparisonImage)
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
			lastDenoisedImage := getLastDenoisedImage(imageQueue)

			if lastDenoisedImage != nil {
				snr = analysis.SignalToNoiseRatio(lastDenoisedImage, comparisonImage)
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
			lastDenoisedImage := getLastDenoisedImage(imageQueue)

			if lastDenoisedImage != nil {
				psnr = analysis.PeakSignalToNoiseRatio(lastDenoisedImage, comparisonImage)
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
			lastDenoisedImage := getLastDenoisedImage(imageQueue)

			if lastDenoisedImage != nil {
				md = analysis.MaxDifference(lastDenoisedImage, comparisonImage)
			} else {
				md = analysis.MaxDifference(img, comparisonImage)
			}

			cmdResult.Description = "Max Difference calculated"
			cmdResult.Result = fmt.Sprintf("Max Difference: %d", md)

		case "histogram":

			outputFileName := fmt.Sprintf("%s_histogram.bmp", originalNameWithoutExt)
			newImg := manipulations.GenerateGraphicalRepresentationOfHistogram(manipulations.CalculateHistogram(img))

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName, IsHistogram: true})

			cmdResult.Description = "Computed Graphical Representation of Histogram"

		case "cmean", "cvariance", "cstdev", "cvarcoi", "casyco", "cflatco", "cvarcoii", "centropy":

			var histogramImg *image.RGBA
			var histogramImgFilename string

			for i := 0; i < len(imageQueue); i++ {
				if imageQueue[i].IsHistogram {
					histogramImg = imageQueue[i].Image
					histogramImgFilename = imageQueue[i].Filename
					break
				}
			}

			var histogram [256]int

			if histogramImg == nil {
				histogram = manipulations.CalculateHistogram(img)
				histogramImgFilename = fmt.Sprintf("%s_histogram.bmp", originalNameWithoutExt)
			} else {
				histogram = manipulations.CalculateHistogram(histogramImg)
			}

			cmdResult.Result, cmdResult.Description = analysis.CalculateHistogramCharacteristic(command.Name, histogram, histogramImgFilename)

		case "hrayleigh":

			gMin := GetOrDefault(command.Args["min"], 0)
			gMax := GetOrDefault(command.Args["max"], 255)
			alpha := GetOrDefault(command.Args["alpha"], 100.0)

			if gMin < 0 || gMax > 255 || gMin >= gMax {
				log.Fatal("gMin and gMax must be in the range [0, 255] with gMin < gMax")
			}

			outputFileName := fmt.Sprintf("%s_rayleigh_min%d_max%d_alpha%.2f.bmp", originalNameWithoutExt, gMin, gMax, alpha)

			newImg := manipulations.EnhanceImageWithRayleigh(img, float64(gMin), float64(gMax), alpha)

			if commands.Includes("histogram") {
				histogramImgAfterTransformation := manipulations.GenerateGraphicalRepresentationOfHistogram(manipulations.CalculateHistogram(newImg))
				histogramFilename := fmt.Sprintf("%s_histogram_after_rayleigh_min%d_max%d_alpha%.2f.bmp", originalNameWithoutExt, gMin, gMax, alpha)

				imageQueue = append(imageQueue, ImageQueueItem{Image: histogramImgAfterTransformation, Filename: histogramFilename, IsHistogram: true})
			}

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

			cmdResult.Description = fmt.Sprintf("Rayleigh transformation applied with gMin: %v, gMax: %v, and alpha: %.3f", gMin, gMax, alpha)

		case "sedgesharp":

			chosenMask := GetOrDefault(command.Args["mask"], "mask1")

			mask, err := manipulations.GetMask(loadedMasks, chosenMask)

			if err != nil {
				log.Fatalf("Error getting mask: %v", err)
			}

			outputFileName := fmt.Sprintf("%s_sharpened_edges_%s.bmp", originalNameWithoutExt, chosenMask)

			newImg := manipulations.ApplyConvolutionUniversal(img, mask)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

		case "okirsf":

			outputFileName := fmt.Sprintf("%s_kirsh_edge_detection.bmp", originalNameWithoutExt)

			newImg := manipulations.ApplyKirshEdgeDetection(img)

			imageQueue = append(imageQueue, ImageQueueItem{Image: newImg, Filename: outputFileName})

		default:
			fmt.Println("Unknown commend")
			return
		}

		cmdResult.Duration = time.Since(startTime)
		commandResults = append(commandResults, cmdResult)

		durationSum += cmdResult.Duration
	}

	for _, imgItem := range imageQueue {
		err = imageio.SaveBmpImage(imgItem.Image, imgItem.Filename)
		if err != nil {
			log.Fatalf("\nError saving file: %v", err)
		} else {
			fmt.Printf("\nImage saved successfully as: %s\n", imgItem.Filename)
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

func getLastDenoisedImage(queue []ImageQueueItem) *image.RGBA {
	for i := len(queue) - 1; i >= 0; i-- {
		if queue[i].Denoised {
			return queue[i].Image
		}
	}
	return nil
}
