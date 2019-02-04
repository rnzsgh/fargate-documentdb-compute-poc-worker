package math

import (
	"testing"
)

func TestMultiply(t *testing.T) {
	t.Run("TestMultiply", func(t *testing.T) {
		_, err := Multiply(
			[][]float32{[]float32{1.0, 2.0, 3.0}, []float32{4.0, 5.0, 6.0}},
			Transpose([][]float32{[]float32{0.5, 0.2, 0.7}, []float32{0.5, 0.8, 0.3}}),
		)
		if err != nil {
			t.Error(err)
		}
	})
}
