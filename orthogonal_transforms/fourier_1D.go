package orthogonal_transforms

import (
	"math"
	"math/cmplx"
)

// DFT computes the Discrete Fourier Transform (DFT) of a signal.
func DFT(signal []complex128) []complex128 {
	N := len(signal)
	result := make([]complex128, N)
	for k := 0; k < N; k++ {
		var sum complex128
		for n := 0; n < N; n++ {
			angle := -2 * math.Pi * float64(k*n) / float64(N)
			sum += signal[n] * cmplx.Exp(complex(0, angle))
		}
		result[k] = sum
	}
	return result
}

// IDFT computes the Inverse Discrete Fourier Transform (IDFT) of a signal.
func IDFT(frequencySignal []complex128) []complex128 {
	N := len(frequencySignal)
	result := make([]complex128, N)
	for n := 0; n < N; n++ {
		var sum complex128
		for k := 0; k < N; k++ {
			angle := 2 * math.Pi * float64(k*n) / float64(N)
			sum += frequencySignal[k] * cmplx.Exp(complex(0, angle))
		}
		result[n] = sum / complex(float64(N), 0)
	}
	return result
}

// FFT computes the Fast Fourier Transform (FFT) using recursion.
func FFT(signal []complex128) []complex128 {
	N := len(signal)
	if N <= 1 {
		return signal
	}
	if N%2 != 0 {
		panic("Signal length must be a power of 2 for FFT.")
	}

	even := FFT(signal[0:N:2])
	odd := FFT(signal[1:N:2])

	result := make([]complex128, N)
	for k := 0; k < N/2; k++ {
		angle := -2 * math.Pi * float64(k) / float64(N)
		t := cmplx.Exp(complex(0, angle)) * odd[k]
		result[k] = even[k] + t
		result[k+N/2] = even[k] - t
	}
	return result
}

// IFFT computes the Inverse Fast Fourier Transform (IFFT).
func IFFT(frequencySignal []complex128) []complex128 {
	N := len(frequencySignal)
	conjugated := make([]complex128, N)
	for i := range frequencySignal {
		conjugated[i] = cmplx.Conj(frequencySignal[i])
	}

	result := FFT(conjugated)
	for i := range result {
		result[i] = cmplx.Conj(result[i]) / complex(float64(N), 0)
	}
	return result
}
