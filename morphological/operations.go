package morphological

func Dilation(image BinaryImage, se StructuringElement) BinaryImage {
	rows := len(image)
	cols := len(image[0])
	output := make(BinaryImage, rows)
	for i := range output {
		output[i] = make([]int, cols)
	}

	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			if image[x][y] == 1 {
				for i := 0; i < len(se.Data); i++ {
					for j := 0; j < len(se.Data[i]); j++ {
						if se.Data[i][j] == 1 {
							newX := x + i - se.OriginX
							newY := y + j - se.OriginY
							if newX >= 0 && newX < rows && newY >= 0 && newY < cols {
								output[newX][newY] = 1
							}
						}
					}
				}
			}
		}
	}

	return output
}

func Erosion(image BinaryImage, se StructuringElement) BinaryImage {
	rows := len(image)
	cols := len(image[0])
	output := make(BinaryImage, rows)
	for i := range output {
		output[i] = make([]int, cols)
	}

	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			if Fits(image, se, x, y) {
				output[x][y] = 1
			}
		}
	}

	return output
}

func Opening(image BinaryImage, se StructuringElement) BinaryImage {
	eroded := Erosion(image, se)
	return Dilation(eroded, se)
}

// Reference: https://www.geeksforgeeks.org/difference-between-opening-and-closing-in-digital-image-processing/
func Closing(image BinaryImage, se StructuringElement) BinaryImage {
	dilated := Dilation(image, se)
	return Erosion(dilated, se)
}

func HitOrMiss(image BinaryImage, se1, se2 StructuringElement) BinaryImage {
	foreground := Erosion(image, se1)
	complement := Complement(image)
	background := Erosion(complement, se2)
	return Intersection(foreground, background)
}

// func Complement(image BinaryImage) BinaryImage {
// 	rows := len(image)
// 	cols := len(image[0])
// 	output := make(BinaryImage, rows)
// 	for i := range output {
// 		output[i] = make([]int, cols)
// 		for j := range output[i] {
// 			// output[i][j] = 1 - image[i][j]
// 			if image[i][j] == 1 {
// 				output[i][j] = 0
// 			} else {
// 				output[i][j] = 1
// 			}
// 		}
// 	}

// 	return output
// }

func Complement(img BinaryImage) BinaryImage {
	for i := range img {
		for j := range img[i] {
			if img[i][j] == 1 {
				img[i][j] = 0
			} else {
				img[i][j] = 1
			}
		}
	}
	return img
}

// func Intersection(img1, img2 BinaryImage) BinaryImage {
// 	rows := len(img1)
// 	cols := len(img1[0])
// 	output := make(BinaryImage, rows)
// 	for i := range output {
// 		output[i] = make([]int, cols)
// 		for j := range output[i] {
// 			// if img1[i][j] == 1 && img2[i][j] == 1 {
// 			// 	output[i][j] = 1
// 			// }
// 			if img1[i][j] == 1 && img2[i][j] == 1 {
// 				output[i][j] = 1
// 			} else {
// 				output[i][j] = 0
// 			}
// 		}
// 	}
// 	return output
// }

func Intersection(img1, img2 BinaryImage) BinaryImage {
	for i := range img1 {
		for j := range img1[i] {
			if img1[i][j] == 1 && img2[i][j] == 0 {
				img1[i][j] = 0
			}
		}
	}
	return img1
}
