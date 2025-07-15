package mappingrules

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Integer | constraints.Float
}

func Abs[T Numeric](value T) T {
	return max(value, -value)
}

// LerpInt linearly interpolates between two integer values using
// a float64 index value
func LerpInt[T constraints.Integer](min, max T, t float64) T {
	t = Clamp(t, 0.0, 1.0)
	return T((1-t)*float64(min) + t*float64(max))
}

func Clamp[T Numeric](value, min, max T) T {
	if value < min {
		value = min
	}
	if value > max {
		value = max
	}
	return value
}
