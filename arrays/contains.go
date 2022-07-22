package arrays

func Contains[T comparable](val T, arr []T) bool {
	for _, x := range arr {
		if x == val {
			return true
		}
	}

	return false
}

func ContainsP[T comparable](val *T, arr []T) bool {
	if val == nil {
		return false
	}

	return Contains(*val, arr)
}
