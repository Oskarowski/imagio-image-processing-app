package orthogonal_transforms

import (
	"math"
	"math/cmplx"
)

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
