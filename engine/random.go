package engine

import "math/rand"

func RandomRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomIntRange(min, max int) int {
	if (min == 0 && max == 0) || max-min <= 0 {
		return 0
	}
	return min + rand.Intn(max-min)
}
