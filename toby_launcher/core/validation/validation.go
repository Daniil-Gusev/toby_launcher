package validation

import (
	"math"
	"strconv"
	"strings"
	"toby_launcher/apperrors"
)

func ParseInt(input string) (int, error) {
	input = strings.TrimSpace(input)
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0, apperrors.New(apperrors.Err, "You must enter an integer.", nil)
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, apperrors.New(apperrors.Err, "The provided value is out of acceptable bounds.", nil)
	}
	return num, nil
}

func IsNumInRange(num, min, max int) (bool, error) {
	if num > math.MaxInt32 || num < math.MinInt32 {
		return false, apperrors.New(apperrors.Err, "The provided value is out of acceptable bounds.", nil)
	}
	if num < min {
		return false, apperrors.New(apperrors.Err, "The provided number must not be less than $min.", map[string]any{
			"min": min,
		})
	}
	if num > max {
		return false, apperrors.New(apperrors.Err, "The provided number must not exceed $max.", map[string]any{
			"max": max,
		})
	}
	return true, nil
}

func ParseIntInRange(input string, min, max int) (int, error) {
	num, err := ParseInt(input)
	if err != nil {
		return 0, err
	}
	if _, err := IsNumInRange(num, min, max); err != nil {
		return 0, err
	}
	return num, nil
}

func ParseOptionalIntInRange(input string, defaultValue, min, max int) (int, error) {
	if input == "" {
		if _, err := IsNumInRange(defaultValue, min, max); err != nil {
			return 0, err
		}
		return defaultValue, nil
	}
	return ParseIntInRange(input, min, max)
}
