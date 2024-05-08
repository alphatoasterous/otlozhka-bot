package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func UnixToTime(unixTime int64, location *time.Location) time.Time {
	return time.Unix(unixTime, 0).In(location)
}

func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// RemoveTrailingComma checks the given string `s` for a trailing comma.
// If `s` ends with a comma, it removes the trailing comma and returns the modified string.
// If there is no trailing comma, it returns the string unchanged.
func RemoveTrailingComma(s string) string {
	if len(s) > 0 && s[len(s)-1] == ',' {
		return s[:len(s)-1]
	}
	return s
}

func GetRandomItemFromStrArray(arr []string) string {
	// I used to roll the dice.
	return arr[rand.Intn(len(arr))]
}
