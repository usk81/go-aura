package common

// InStrings checks if a value exists in an string slice
func InStrings(k string, ss []string) bool {
	for _, s := range ss {
		if s == k {
			return true
		}
	}
	return false
}
