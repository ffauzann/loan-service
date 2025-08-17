package str

import (
	"fmt"
	"strings"
)

func PaddingNumberZero(width int, value string) string {
	return fmt.Sprintf("%0*s", width, value)
}

func PhoneWithCountryCode(value string, defaultCountryCode string, removeSeparator bool) string {
	if removeSeparator {
		value = RemoveSeparator(value)
	}

	switch value[:1] {
	case "0":
		return defaultCountryCode + value[1:]
	case "+":
		return value
	default:
		return "+" + value
	}
}

func RemoveSeparator(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "_", "")

	return s
}
