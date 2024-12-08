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

	image := morphological.BinaryImage{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0},
		{0, 1, 1, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 1, 1, 1, 0},
		{0, 0, 0, 1, 1, 1, 1, 0},
		{0, 0, 0, 1, 0, 1, 0, 0},
		{0, 0, 1, 1, 0, 1, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
	}

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
	fmt.Println("Dilated:")
	for _, row := range dilated {
		fmt.Println(row)
	}

	eroded := morphological.Erosion(image, se)
	fmt.Println("Eroded:")
	for _, row := range eroded {
		fmt.Println(row)
	}

	opened := morphological.Opening(image, se)
	fmt.Println("Opened:")
	for _, row := range opened {
		fmt.Println(row)
	}

	closed := morphological.Closing(image, se)
	fmt.Println("Closed:")
	for _, row := range closed {
		fmt.Println(row)
	}

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
	fmt.Println("Hit or Miss:")
	for _, row := range resultHoS {
		fmt.Println(row)
	}

	fmt.Println()

	A := morphological.BinaryImage{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	fmt.Println("Before Thinning Image:")
	for _, rowA := range A {
		fmt.Println(rowA)
	}

	thinnedImage := morphological.Thinning(A, morphological.StructuralElements)
	fmt.Println("\nThinned Image:")
	for _, row := range thinnedImage {
		fmt.Println(row)
	}

	seeds := []morphological.Point{{X: 100, Y: 300}}

	// hot to print the path?
	absolutePath, _ := filepath.Abs("imgs/mandrilc.bmp")
	img, err := imageio.OpenBmpImage(absolutePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	fmt.Println("Img opened...")

	_, segmentedImg := morphological.RegionGrowing(img, seeds, 0, 40)
	fmt.Println("Saving segmented image...")
	imageio.SaveBmpImage(segmentedImg, "segmented_100_300_40.bmp")
	fmt.Println("Segmented image SAVED!")

	return

	if len(os.Args) > 1 {
		cmd.RunAsCliApp()
	} else {
		gui.RunAsTUIApp()
	}
}
