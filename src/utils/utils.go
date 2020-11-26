package utils

func StringIsContained(slice []string, queried string) bool {
	// Search in an array of strings a particular array
	for _, a := range slice {
		if a == queried {
			return true
		}
	}
	return false
}
