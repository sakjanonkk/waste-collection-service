package utils

import (
	"errors"
	"strconv"
)

// StringToUint แปลง string → uint
func StringToUint(s string, out *uint) error {
	if s == "" {
		return errors.New("empty string")
	}

	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}

	*out = uint(val)
	return nil
}
