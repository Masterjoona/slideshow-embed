package util

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func UrlDecodeString(input string) string {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return input
	}
	return decoded
}

func UrlDecodeStrings(input []string) []string {
	for i, str := range input {
		input[i] = UrlDecodeString(str)
	}
	return input
}

func SplitURLAndIndex(URL string) (string, string, bool) {
	lastInd := strings.LastIndex(URL, "/")
	index := URL[lastInd+1:]
	if index == "" {
		index = "1"
	}
	sound := strings.HasSuffix(index, "s")
	if sound {
		index = strings.Replace(index, "s", "", 1)
	}
	return URL[:lastInd], index, sound
}

func Ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

const million = 1000000
const thousand = 1000

func FormatLargeNumbers(numberString string) string {
	number, err := strconv.Atoi(numberString)
	if err != nil {
		return "0"
	}
	switch {
	case number >= million:
		return fmt.Sprintf("%.1fM", float64(number)/million)
	case number >= thousand:
		return fmt.Sprintf("%.1fK", float64(number)/thousand)
	default:
		return fmt.Sprintf("%d", number)
	}
}
