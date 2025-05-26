package main

import (
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key := make([]byte, 32)

	mgm, err := mgmgo.New(key)
	if err != nil {
		log.Fatalf("Error during encryptor creation: %s", err)
	}

	additionalData := []byte("your-additional-data")
	plaintext := []byte("your-message")
	nonce := make([]byte, mgm.NonceSize())

	result := mgm.Seal(plaintext[:0], nonce, plaintext, additionalData)

	fmt.Printf("Ciphertext: %x\n", result[:len(plaintext)])
	fmt.Printf("MAC: %x\n", result[len(plaintext):])
}
