package morphological

func N(A BinaryImage, B StructuringElement) BinaryImage {
	eroded := Erosion(A, B)
	complement := Complement(eroded)

	return Intersection(A, complement)
}

func Equal(img1, img2 BinaryImage) bool {
	if len(img1) != len(img2) || len(img1[0]) != len(img2[0]) {
		return false
	}

	for i := range img1 {
		for j := range img1[i] {
			if img1[i][j] != img2[i][j] {
				return false
			}
		}
	}

	return true
}

func NRepeated(A BinaryImage, SEs []StructuringElement) BinaryImage {
	prev := A
	for {
		current := A
		for _, B := range SEs {
			current = N(current, B)
		}

		if Equal(prev, current) {
			break
		}

		prev = current
	}

	return prev
}
