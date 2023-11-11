package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func Int64SliceToString(slice []int64) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprint(v)
	}
	return strings.Join(strSlice, ",")
}

func Int64SliceToStringLimit5(slice []int64) string {
	var strSlice []string
	for i, v := range slice {
		if i >= 5 {
			break
		}
		strSlice = append(strSlice, strconv.FormatInt(v, 10))
	}
	return strings.Join(strSlice, ",")
}

func Float64SliceToStringLimit5(slice []float64) string {
	var strSlice []string
	for i, v := range slice {
		if i >= 5 {
			break
		}
		strSlice = append(strSlice, strconv.FormatFloat(v, 'f', -1, 64))
	}
	return strings.Join(strSlice, ",")
}
