package sportmonks

import (
	"strconv"
	"strings"
)

//IntSliceToSepString takes a int slice and generates a string of values separated by specified separator
func IntSliceToSepString(slice []int, sep string) string {
	if len(slice) == 0 {
		return ""
	}

	s := make([]string, len(slice))
	for i, v := range slice {
		s[i] = strconv.Itoa(v)
	}
	return strings.Join(s, sep)
}
