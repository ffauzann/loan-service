package util

import "math/rand"

var (
	alphabets = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers   = []rune("1234567890")
)

func RandomAlphaSequence(seed int64, length uint8) string {
	r := rand.New(rand.NewSource(seed))

	b := make([]rune, length)
	for i := range b {
		b[i] = alphabets[r.Intn(len(alphabets))]
	}

	return string(b)
}

func RandomNumericSequence(seed int64, length uint8) string {
	r := rand.New(rand.NewSource(seed))

	b := make([]rune, length)
	for i := range b {
		b[i] = numbers[r.Intn(len(numbers))]
	}

	return string(b)
}
