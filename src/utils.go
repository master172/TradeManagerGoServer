package main

import (
	"crypto/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

const codeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const codeLength = 12

// GenCode returns a cryptographically random string of length codeLength.
func GenCode() string {
	b := make([]byte, codeLength)
	_, err := rand.Read(b)
	if err != nil {
		// fallback, extremely unlikely
		for i := range b {
			b[i] = byte(codeChars[0])
		}
		return string(b)
	}
	for i := range codeLength {
		b[i] = codeChars[int(b[i])%len(codeChars)]
	}
	return string(b)
}

var upgrader = websocket.Upgrader{
	// In production you should validate the Origin header
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict allowed origins in production
		return true
	},
}
