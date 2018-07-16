package gather

// DoesSliceContainString - checks a string slice to see if it contains the provided string.
func DoesSliceContainString(container []string, value string) bool {
	for _, s := range container {
		if s == value {
			return true
		}
	}
	return false
}
