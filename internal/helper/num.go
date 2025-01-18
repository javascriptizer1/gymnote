package helper

import (
	"strconv"
	"strings"
)

func ParseFloat32(s string) (float32, error) {
	val, err := strconv.ParseFloat(strings.TrimSpace(s), 32)
	if err != nil {
		return 0, err
	}
	return float32(val), nil
}

func ParseUint8(s string) (uint8, error) {
	val, err := strconv.ParseUint(strings.TrimSpace(s), 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(val), nil
}
