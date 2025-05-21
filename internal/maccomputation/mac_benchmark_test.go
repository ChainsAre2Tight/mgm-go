package maccomputation_test

import (
	"context"
	"encoding/hex"
	"testing"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/maccomputation"
)

func BenchmarkMAC(b *testing.B) {
	key, err := hex.DecodeString("8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF")
	if err != nil {
		b.Fatalf("Error during key decoding: %s", err)
	}
	k, err := kuznechikgo.Schedule(key)
	if err != nil {
		b.Fatalf("Error during keyschedule: %s", err)
	}
	keys := kuznechikgo.KeysToUints(k)
	nonce := bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88")
	authenticatedData := []*bitstrings.BitString128{
		bitstrings.FromGOSTString("02 02 02 02 02 02 02 02 01 01 01 01 01 01 01 01"),
		bitstrings.FromGOSTString("04 04 04 04 04 04 04 04 03 03 03 03 03 03 03 03"),
		bitstrings.FromGOSTString("EA 05 05 05 05 05 05 05 05 00 00 00 00 00 00 00"),
	}
	ciphertext := []*bitstrings.BitString128{
		bitstrings.FromGOSTString("A9 75 7B 81 47 95 6E 90 55 B8 A3 3D E8 9F 42 FC"),
		bitstrings.FromGOSTString("80 75 D2 21 2B F9 FD 5B D3 F7 06 9A AD C1 6b 39"),
		bitstrings.FromGOSTString("49 7A B1 59 15 A6 BA 85 93 6B 5D 0E A9 F6 85 1C"),
		bitstrings.FromGOSTString("C6 0C 14 D4 D3 F8 83 D0 AB 94 42 06 95 C7 6D EB"),
		bitstrings.FromGOSTString("2C 75 52 00 00 00 00 00 00 00 00 00 00 00 00 00"),
	}
	lengthAuth := uint64(0x0000000000000148)
	lengthPlain := uint64(0x0000000000000218)
	ctx := context.Background()
	for b.Loop() {
		_, err := maccomputation.Compute(
			keys, nonce, authenticatedData, ciphertext, lengthAuth, lengthPlain, ctx,
		)
		if err != nil {
			b.Fatalf("Error during maccompute: %s", err)
		}
	}
}
