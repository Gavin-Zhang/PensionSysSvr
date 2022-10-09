package unmarshal

import (
	"strconv"
	"strings"
)

// StringToUint string to uint
func StringToUint(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return i, err
	}
	return 0, err
}

// StringToInt string to int
func StringToInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i, err
	}
	return 0, err
}

// StringToFloat string to float
func StringToFloat(s string) (float64, error) {
	i, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return i, err
	}
	return 0, err
}

// StringToBool string to bool
func StringToBool(s string) (bool, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return false, err
	}
	if i == 0 {
		return false, nil
	}
	return true, nil
}

//
func StringDelEnter(s string) string {
	index := strings.IndexAny(s, "\n")
	if index != -1 {
		return s[:index] + s[index+1:]
	}
	return s
}

// StringToSlice string to []string
func StringToSlice(s string, split rune) ([]string, error) {
	canSplit := func(c rune) bool { return c == split }
	Slice := strings.FieldsFunc(s, canSplit)
	return Slice, nil
}

// StringToUint64Slice string to []uint64, string format [string, string]
func StringToUint64Slice(s string) ([]uint64, error) {
	s = strings.Trim(s, "[")
	s = strings.Trim(s, "]")

	params, _ := StringToSlice(s, ',')
	backslice := make([]uint64, 0)
	for i := 0; i < len(params); i++ {
		param, err := StringToUint(params[i])
		if err != nil {
			return backslice, err
		}
		backslice = append(backslice, param)
	}
	return backslice, nil
}
