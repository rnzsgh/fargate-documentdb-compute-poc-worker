package math

import (
	"errors"
)

// Source from corrected version on:  https://gist.github.com/n1try/c5082f0f1db7f4abcb6d995dc275fe7f
func Transpose(x [][]float32) [][]float32 {
	out := make([][]float32, len(x[0]))
	for i := 0; i < len(x); i += 1 {
		for j := 0; j < len(x[0]); j += 1 {
			out[j] = append(out[j], x[i][j])
		}
	}
	return out
}

// Source from corrected version on:  https://gist.github.com/n1try/c5082f0f1db7f4abcb6d995dc275fe7f
func Multiply(x, y [][]float32) ([][]float32, error) {
	if len(x[0]) != len(y) {
		return nil, errors.New("Can't do matrix multiplication - entry lengths are not the same")
	}

	out := make([][]float32, len(x))
	for i := 0; i < len(x); i++ {
		out[i] = make([]float32, len(y[0]))
		for j := 0; j < len(y[0]); j++ {
			for k := 0; k < len(y); k++ {
				out[i][j] += x[i][k] * y[k][j]
			}
		}
	}
	return out, nil
}
