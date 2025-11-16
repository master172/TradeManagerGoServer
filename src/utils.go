package main

import (
	"math/rand"
)

const codeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const codeLength = 12

func GenCode() string {
	result := make([]byte, codeLength)

	for i := range codeLength {
		result[i] = codeChars[rand.Intn(len(codeChars))]
	}

	return string(result)
}
