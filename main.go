package main

import (
	"fmt"
	"image-processing/cmd"
	"image-processing/cmd/gui"
	"image-processing/imageio"
	"image-processing/morphological"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// if len(os.Args) > 1 && os.Args[1] == "--help" {
	// 	cmd.PrintHelp()
	// 	return
	// }

	logFile, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	absoluteBinaryImgPath, _ := filepath.Abs("imgs/lenabw.bmp")
	img, err := imageio.LoadMonochromeBMP(absoluteBinaryImgPath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	image := morphological.ConvertIntoBinaryImage(img)

	imageio.SaveBmpImage(morphological.ConvertIntoImage(image), "converted_to_binary_and_into_image.bmp")

	se := morphological.StructuringElement{
		Data: [][]int{
			{0, 1, 0},
			{1, 1, 1},
			{0, 1, 0},
		},
		OriginX: 1,
		OriginY: 1,
	}

	dilated := morphological.Dilation(image, se)
	imageio.SaveBmpImage(morphological.ConvertIntoImage(dilated), "dilated.bmp")
	fmt.Println("Dilated:")

	eroded := morphological.Erosion(image, se)
	imageio.SaveBmpImage(morphological.ConvertIntoImage(eroded), "eroded.bmp")

	opened := morphological.Opening(image, se)
	imageio.SaveBmpImage(morphological.ConvertIntoImage(opened), "opened.bmp")

	closed := morphological.Closing(image, se)
	imageio.SaveBmpImage(morphological.ConvertIntoImage(closed), "closed.bmp")

	foregroundSE := morphological.StructuringElement{
		Data: [][]int{
			{0, 1, 0},
			{0, 1, 0},
			{1, 1, 0},
		},
		OriginX: 2, OriginY: 1,
	}

	backgroundSE := morphological.StructuringElement{
		Data: [][]int{
			{0, 0, 0},
			{1, 0, 0},
			{0, 0, 0},
		},
		OriginX: 2, OriginY: 1,
	}
	resultHoS := morphological.HitOrMiss(image, foregroundSE, backgroundSE)
	imageio.SaveBmpImage(morphological.ConvertIntoImage(resultHoS), "resultHoS.bmp")

	fmt.Println()

	thinnedImage := morphological.Thinning(image, morphological.StructuralElements)
	imageio.SaveBmpImage(morphological.ConvertIntoImage(thinnedImage), "thinned.bmp")

	absolutePath, _ := filepath.Abs("imgs/lenac.bmp")
	imgg, err := imageio.OpenBmpImage(absolutePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	seeds := []morphological.Point{{X: 320, Y: 450}, {X: 500, Y: 40}, {X: 50, Y: 250}, {X: 200, Y: 236}}
	_, segmentedImg := morphological.RegionGrowing(imgg, seeds, morphological.Chebyshev, 35)

	fmt.Println("Saving segmented image...")
	imageio.SaveBmpImage(segmentedImg, "segmented_320_450_30_2.bmp")
	fmt.Println("Segmented image SAVED!")

	return

	if len(os.Args) > 1 {
		cmd.RunAsCliApp()
	} else {
		gui.RunAsTUIApp()
	}
}
