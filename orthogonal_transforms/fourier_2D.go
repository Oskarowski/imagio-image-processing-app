package orthogonal_transforms

func SlowDFT2D(input [][]complex128, inverse bool) [][]complex128 {
	n := len(input)
	m := len(input[0])
	output := make([][]complex128, n)
	for i := range output {
		output[i] = make([]complex128, m)
	}

	for i := 0; i < n; i++ {
		output[i] = SlowDFT1D(input[i], inverse)
	}

	for j := 0; j < m; j++ {
		column := make([]complex128, n)
		for i := 0; i < n; i++ {
			column[i] = output[i][j]
		}
		column = SlowDFT1D(column, inverse)
		for i := 0; i < n; i++ {
			output[i][j] = column[i]
		}
	}

	return output
}

// Transpose returns the transpose of a given 2D complex matrix.
// It swaps the rows and columns of the input matrix, effectively
// flipping it over its diagonal. The resulting matrix has dimensions
// where the number of rows and columns are interchanged.

func transpose(matrix [][]complex128) [][]complex128 {
	n := len(matrix)
	m := len(matrix[0])
	transpose := make([][]complex128, m)
	for i := 0; i < m; i++ {
		transpose[i] = make([]complex128, n)
		for j := 0; j < n; j++ {
			transpose[i][j] = matrix[j][i]
		}
	}
	return transpose
}

func FFT2D(input [][]complex128, inverse bool) [][]complex128 {
	n := len(input)
	m := len(input[0])

	// Step 1: Apply 1D FFT to each row
	for i := 0; i < n; i++ {
		input[i] = FFT1D(input[i], inverse)
	}

	// Step 2: Transpose the matrix
	transposed := transpose(input)

	// Step 3: Apply 1D FFT to each column (now rows of the transposed matrix)
	for i := 0; i < m; i++ {
		transposed[i] = FFT1D(transposed[i], inverse)
	}

	// Step 4: Transpose back to original orientation
	return transpose(transposed)
}
