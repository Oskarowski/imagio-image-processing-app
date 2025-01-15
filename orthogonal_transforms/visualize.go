package orthogonal_transforms

import (
	"math/cmplx"
)

func VisualizeSpectrum(input [][]complex128) [][]float64 {
	N, M := len(input), len(input[0])
	output := make([][]float64, N)
	for i := range output {
		output[i] = make([]float64, M)
		for j := range output[i] {
			output[i][j] = cmplx.Abs(input[i][j])
		}
	}
	return output
}
