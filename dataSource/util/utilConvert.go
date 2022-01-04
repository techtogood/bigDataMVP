package util

import "strconv"

func ConvertFloat64(s string) (float float64) {
	float, _ = strconv.ParseFloat(s, 64)
	return
}


func ConvertString(i int) (s string) {
	s = strconv.Itoa(i)
	return
}
