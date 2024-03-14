package utils

import (
	"math/rand"
	"strconv"
)

func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func RemoveTrailingComma(s string) string {
	// TODO: Ewwww
	if len(s) > 0 && s[len(s)-1] == ',' {
		return s[:len(s)-1]
	}
	return s
}

func GetRandomItemFromStrArray(arr []string) string {
	// I used to roll the dice.
	return arr[rand.Intn(len(arr))]
}
