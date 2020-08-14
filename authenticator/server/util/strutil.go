package util

import "strings"

// RemoveDuplicatesStable removes duplicate and empty elements from a slice of
// strings, preserving order (and case) of the original slice.
// In all cases, strings are compared after trimming whitespace
// If caseInsensitive, strings will be compared after ToLower().
func RemoveDuplicatesStable(items []string, caseInsensitive bool) []string {
	itemsMap := make(map[string]bool, len(items))
	deduplicated := make([]string, 0, len(items))

	for _, item := range items {
		key := strings.TrimSpace(item)
		if caseInsensitive {
			key = strings.ToLower(key)
		}
		if key == "" || itemsMap[key] {
			continue
		}
		itemsMap[key] = true
		deduplicated = append(deduplicated, item)
	}
	return deduplicated
}

// StrListDelete removes the first occurrence of the given item from the slice
// of strings if the item exists.
// If caseInsensitive, strings will be compared after ToLower().
func StrListDelete(s []string, d string, caseInsensitive bool) []string {
	if s == nil {
		return s
	}
	if caseInsensitive {
		d = strings.ToLower(d)
	}
	for index, element := range s {
		if caseInsensitive {
			element = strings.ToLower(element)
		}
		if element == d {
			return append(s[:index], s[index+1:]...)
		}
	}
	return s
}
