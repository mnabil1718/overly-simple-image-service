package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateCSRFKey() string {
	key := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(key); err != nil {
		fmt.Println("Error generating CSRF key:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(key)
}

func main() {
	key := generateCSRFKey()
	fmt.Println("Generated CSRF Key:", key)
}
