package main

import (
	"encoding/hex"
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {
	nonceGenerator := mgmgo.NewNonceGenerator()

	encryptor := mgmgo.NewEncryptor(nonceGenerator)

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key, err := hex.DecodeString("8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF")
	if err != nil {
		log.Fatalf("Error during key decoding: %s", err)
	}
	associatedData := []byte("your-associated-data")
	plaintext := []byte("your-message")

	nonce, ciphertext, mac, err := encryptor.Encrypt(key, associatedData, plaintext)
	if err != nil {
		log.Fatalf("Encryption failed: %s", err)
	}

	fmt.Printf("Nonce: %x\n", nonce)
	fmt.Printf("Ciphertext: %x\n", ciphertext)
	fmt.Printf("MAC: %x\n", mac)
}
