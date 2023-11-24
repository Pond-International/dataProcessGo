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

func StringSliceToIntSlice(strs string) ([]int64, error) {
	strSlice := strings.Split(strs, ",")
	var intSlice []int64
	for _, str := range strSlice {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		intSlice = append(intSlice, num)
	}
	return intSlice, nil
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
