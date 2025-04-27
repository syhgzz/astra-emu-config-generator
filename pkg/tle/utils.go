package tle

import (
	"fmt"
	"strconv"
	"strings"
)

// CalculateChecksum calculates the sum of all digits in the line modulo 10
func calculateChecksum(line string) int {
	sum := 0
	for _, c := range line {
		if c >= '0' && c <= '9' {
			sum += int(c - '0')
		}
		if c == '-' {
			sum += 1
		}
	}
	return sum % 10
}

// convertToSpecialExponential converts a float to a string with a special exponent format
// The output format is like exponential notation, but the exponent sign `e` or `E` is not showing.
// Besides, the output assume there is a `dot` before the exponent.
// 0.000021815, 7 -> 21815-4
// 0.00020013815, 10 -> 20013815-3
func ConvertToSpecialExponential(f float64, digits int) string {
	str := fmt.Sprintf("%.*e", digits, f)

	parts := strings.Split(str, "e")
	if len(parts) != 2 {
		return str
	}

	mantissa := parts[0]
	exp, _ := strconv.Atoi(parts[1])

	mantissa_parts := strings.Split(mantissa, ".")
	mantissa_int, _ := strconv.Atoi(mantissa_parts[0])
	if mantissa_int != 0 {
		exp = exp + 1
	}

	mantissa = strings.Replace(mantissa, ".", "", 1)

	restDigits := digits - len(strconv.Itoa(exp))
	if restDigits <= 0 {
		return fmt.Sprintf("%s%+d", mantissa, exp)
	}

	if len(mantissa) > restDigits {
		mantissa = mantissa[:restDigits]
	} else if len(mantissa) < restDigits {
		mantissa = mantissa + strings.Repeat("0", restDigits-len(mantissa))
	}

	return fmt.Sprintf("%s%+d", mantissa, exp)
}

func GetSign(f float64) string {
	if f >= 0 {
		return "+"
	}
	return "-"
}

// ConvertToDotfree converts a float to a string without the integer part
// 0.00001478, 7 -> 0001478
func ConvertToDotfree(f float64, digits int) string {
	str := fmt.Sprintf("%.*f", digits, f)
	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return str
	}
	return parts[1]
}
