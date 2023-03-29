package slice

func Containers[T comparable](find T, slice []T) bool {
	for _, item := range slice {
		if find == item {
			return true
		}
	}

	return false
}
