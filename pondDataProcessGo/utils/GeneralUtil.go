package utils

import (
	"fmt"
	"strings"
)

func Int64SliceToString(slice []int64) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}
	return strings.Join(strSlice, ",")
}
