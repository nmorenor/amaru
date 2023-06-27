package engine

func Clamp(value, min, max float64) (float64, bool) {
	if value < min {
		return min, true
	}
	if value > max {
		return max, true
	}
	return value, (value <= min || value >= max)
}
