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

	fmt.Print("Enter nonce (32 hex symbols): ")
	text, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	text = text[:len(text)-1]
	nonce, err := hex.DecodeString(text)
	if err != nil {
		log.Fatalf("Error during nonce decoding: %s", err)
	} else if l := len(nonce); l != 16 {
		log.Fatalf("Invalid nonce length: expected: 16, got: %d", l)
	}

	fmt.Print("Enter additional data: ")
	text, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	text = text[:len(text)-1]
	additionalData := []byte(text)

	fmt.Print("Enter plaintext: ")
	text, err = reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	text = text[:len(text)-1]
	plaintext := []byte(text)

	result := mgm.Seal(plaintext[:0], nonce, plaintext, additionalData)

	fmt.Printf("Ciphertext (hex): %x\n", result[:len(plaintext)])
	fmt.Printf("MAC (hex): %x\n", result[len(plaintext):])
}
