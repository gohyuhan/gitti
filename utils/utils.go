package utils

// universal utils that can be used by any package

// Contains is a generic helper function to check for the existence of an item in a slice.
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
