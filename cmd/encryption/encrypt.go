package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter key (64 hex symbols): ")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	text = text[:len(text)-1]
	key, err := hex.DecodeString(text)
	if err != nil {
		log.Fatalf("Error during key decoding: %s", err)
	}

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
