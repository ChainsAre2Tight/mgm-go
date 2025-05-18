package main

import (
	"encoding/hex"
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {
	decryptor := mgmgo.NewDecryptor()

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key, err := hex.DecodeString("8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF")
	if err != nil {
		log.Fatalf("Error during key decoding: %s", err)
	}
	nonce, err := hex.DecodeString("00000000000000000000000000000001")
	if err != nil {
		log.Fatalf("Error during nonce decoding: %s", err)
	}
	associatedData := []byte("your-associated-data")
	ciphertext, err := hex.DecodeString("5c44d197d9aa123feb46d896")
	if err != nil {
		log.Fatalf("Error during ciphertext decoding: %s", err)
	}
	mac, err := hex.DecodeString("5f11114a5ee24fdd5085d6ca11a249fe")
	if err != nil {
		log.Fatalf("Error during mac decoding: %s", err)
	}

	// if MAC authentication fails, an ErrMACsDiffer will be returned
	plaintext, err := decryptor.Decrypt(key, nonce, associatedData, ciphertext, mac)
	if err != nil {
		log.Fatalf("Decryption failed: %s", err)
	}

	fmt.Printf("Plaintext: %s\n", string(plaintext))
}
