package orthogonal_transforms

import (
	"math"
	"math/cmplx"
)

// 2D DFT computes the two-dimensional Discrete Fourier Transform of a 2D matrix.
func DFT2D(input [][]complex128) [][]complex128 {
	N := len(input)    // Number of rows
	M := len(input[0]) // Number of columns

	output := make([][]complex128, N)
	for p := 0; p < N; p++ {
		output[p] = make([]complex128, M)
		for q := 0; q < M; q++ {
			var sum complex128
			for n := 0; n < N; n++ {
				for m := 0; m < M; m++ {
					angle := -2 * math.Pi * (float64(p*n)/float64(N) + float64(q*m)/float64(M))
					sum += input[n][m] * cmplx.Exp(complex(0, angle))
				}
			}
			output[p][q] = sum / complex(math.Sqrt(float64(N*M)), 0)
		}
	}
	return output
}

// 2D IDFT computes the two-dimensional Inverse Discrete Fourier Transform of a 2D matrix.
func IDFT2D(input [][]complex128) [][]complex128 {
	N := len(input)    // Number of rows
	M := len(input[0]) // Number of columns

	output := make([][]complex128, N)
	for n := 0; n < N; n++ {
		output[n] = make([]complex128, M)
		for m := 0; m < M; m++ {
			var sum complex128
			for p := 0; p < N; p++ {
				for q := 0; q < M; q++ {
					angle := 2 * math.Pi * (float64(p*n)/float64(N) + float64(q*m)/float64(M))
					sum += input[p][q] * cmplx.Exp(complex(0, angle))
				}
			}
			output[n][m] = sum / complex(math.Sqrt(float64(N*M)), 0)
		}
	}
	return output
}

// Helper function: 1D FFT
func FFT1D(input []complex128, inverse bool) []complex128 {
	n := len(input)
	if n <= 1 {
		return input
	}

	// Divide: Separate input into even and odd indices
	even := make([]complex128, n/2)
	odd := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		even[i] = input[i*2]
		odd[i] = input[i*2+1]
	}

	// Recursively compute FFT for both halves
	evenFFT := FFT1D(even, inverse)
	oddFFT := FFT1D(odd, inverse)

	// Combine: Apply the FFT butterfly computation
	output := make([]complex128, n)
	angle := 2 * math.Pi / float64(n)
	if inverse {
		angle = -angle
	}
	wn := cmplx.Exp(complex(0, angle))
	w := complex(1, 0)

	for i := 0; i < n/2; i++ {
		output[i] = evenFFT[i] + w*oddFFT[i]
		output[i+n/2] = evenFFT[i] - w*oddFFT[i]
		if inverse {
			output[i] /= 2
			output[i+n/2] /= 2
		}
		w *= wn
	}
	return output
}

// Helper function: Transpose a 2D slice
func Transpose(matrix [][]complex128) [][]complex128 {
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

// 2D FFT
func FFT2D(input [][]complex128, inverse bool) [][]complex128 {
	n := len(input)
	m := len(input[0])

	// Step 1: Apply 1D FFT to each row
	for i := 0; i < n; i++ {
		input[i] = FFT1D(input[i], inverse)
	}

	// Step 2: Transpose the matrix
	transposed := Transpose(input)

	// Step 3: Apply 1D FFT to each column (now rows of the transposed matrix)
	for i := 0; i < m; i++ {
		transposed[i] = FFT1D(transposed[i], inverse)
	}

	// Step 4: Transpose back to original orientation
	return Transpose(transposed)
}

// SwapQuadrants centers the DC component of a 2D FFT result
func SwapQuadrants(matrix [][]complex128) [][]complex128 {
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
