package main

import (
	"encoding/hex"
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key := make([]byte, 32)
	nonce := make([]byte, 16)
	additionalData := []byte("your-associated-data")
	ciphertext, _ := hex.DecodeString("705007f92ecfb7cffaf6f009")
	mac, _ := hex.DecodeString("df70a7a01caf7e134f9d6613df9c06c2")

	mgm, err := mgmgo.New(key)
	if err != nil {
		log.Fatalf("Error during encryptor creation: %s", err)
	}

	result, err := mgm.Open(ciphertext[:0], nonce, append(ciphertext, mac...), additionalData)

	if err != nil {
		log.Fatalf("Error: MACs are different")
	}

	fmt.Printf("Plaintext: %s\n", string(result))
}
