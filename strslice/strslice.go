package strslice

// Contains checks string slice contains a string
func Contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// Append adds a string to slice
func Append(list []string, s string) []string {
	return append(list, s)
}

// Remove removes an item from stringlist
func Remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

