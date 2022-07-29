package arrays

func Index[T comparable](val T, arr []T) int {
	for i, x := range arr {
		if x == val {
			return i
		}
	}

	return -1
}

func IndexP[T comparable](val *T, arr []T) int {
	if val == nil {
		return -1
	}

	return Index(*val, arr)
}
