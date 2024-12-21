package utility

// Abs returns absolute value of an integer
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
