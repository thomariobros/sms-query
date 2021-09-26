package util

// Min min between int
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}
