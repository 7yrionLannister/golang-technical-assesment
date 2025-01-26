package util

import (
	"strconv"
	"strings"
)

func String2UintSlice(str string) ([]uint, error) {
	var result []uint
	for _, s := range strings.Split(str, ",") {
		if s == "" {
			continue
		}
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		result = append(result, uint(num))
	}
	return result, nil
}
