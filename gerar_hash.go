package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	body, err := os.ReadFile("payload.json")
	if err != nil {
		panic(err)
	}

	secret := "meu-secret"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)

	hash := hex.EncodeToString(mac.Sum(nil))
	fmt.Println("sha256=" + hash)
}
