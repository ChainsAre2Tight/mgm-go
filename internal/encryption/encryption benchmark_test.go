package encryption_test

import (
	"context"
	"encoding/hex"
	"testing"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/encryption"
)

func BenchmarkEncryption(b *testing.B) {
	key, err := hex.DecodeString("8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF")
	if err != nil {
		b.Fatalf("Error during key decoding: %s", err)
	}
	nonce := bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88")
	plaintext := []*bitstrings.BitString128{
		bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
		bitstrings.FromGOSTString("00 11 22 33 44 55 66 77 88 99 AA BB CC EE FF 0A"),
		bitstrings.FromGOSTString("11 22 33 44 55 66 77 88 99 AA BB CC EE FF 0A 00"),
		bitstrings.FromGOSTString("22 33 44 55 66 77 88 99 AA BB CC EE FF 0A 00 11"),
	}
	keys, err := kuznechikgo.Schedule(key)
	if err != nil {
		b.Fatalf("Error during keyschedule: %s", err)
	}
	ctx := context.Background()
	for b.Loop() {
		_, err := encryption.Encrypt(plaintext, keys, nonce, ctx)
		if err != nil {
			b.Fatalf("Error during encryption: %s", err)
		}
	}
}
