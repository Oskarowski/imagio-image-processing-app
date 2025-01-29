package orthogonal_transforms

// QuadrantsSwap performs a quadrant swap (also known as a centering) on the
// given 2D complex matrix.
//
// The quadrant swap is a bijective mapping of the matrix elements, where
// the top-left quadrant is swapped with the bottom-right quadrant, and
// the top-right quadrant is swapped with the bottom-left quadrant.
//
// The effect of this operation is to "center" the 2D DFT/IDFT, which is
// useful when performing filtering
// in the frequency domain.
func QuadrantsSwap(matrix [][]complex128) [][]complex128 {
	n := len(matrix)
	m := len(matrix[0])

	centered := make([][]complex128, n)
	for i := range centered {
		centered[i] = make([]complex128, m)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			newI := (i + n/2) % n
			newJ := (j + m/2) % m
			centered[newI][newJ] = matrix[i][j]
		}
	}

	return centered
}
