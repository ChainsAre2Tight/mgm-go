package encryption_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/encryption"
)

func TestDecrypt(t *testing.T) {
	tt := []struct {
		key        string
		nonce      *bitstrings.BitString128
		ciphertext []*bitstrings.BitString128
		plaintext  []*bitstrings.BitString128
	}{
		{
			key:   "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			nonce: bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
			plaintext: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
				bitstrings.FromGOSTString("00 11 22 33 44 55 66 77 88 99 AA BB CC EE FF 0A"),
				bitstrings.FromGOSTString("11 22 33 44 55 66 77 88 99 AA BB CC EE FF 0A 00"),
				bitstrings.FromGOSTString("22 33 44 55 66 77 88 99 AA BB CC EE FF 0A 00 11"),
			},
			ciphertext: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("A9 75 7B 81 47 95 6E 90 55 B8 A3 3D E8 9F 42 FC"),
				bitstrings.FromGOSTString("80 75 D2 21 2B F9 FD 5B D3 F7 06 9A AD C1 6b 39"),
				bitstrings.FromGOSTString("49 7A B1 59 15 A6 BA 85 93 6B 5D 0E A9 F6 85 1C"),
				bitstrings.FromGOSTString("C6 0C 14 D4 D3 F8 83 D0 AB 94 42 06 95 C7 6D EB"),
			},
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s, %s | %v -> %v", td.key, td.nonce, td.ciphertext, td.plaintext),
			func(t *testing.T) {
				k, err := hex.DecodeString(td.key)
				if err != nil {
					t.Fatalf("Error during key decoding: %s", err)
				}
				keys, err := kuznechikgo.Schedule(k)
				if err != nil {
					t.Fatalf("error during keyschedule: %s", err)
				}
				if res, err := encryption.Decrypt(td.ciphertext, keys, td.nonce, context.Background()); err != nil {
					t.Fatalf("error during decryption: %s", err)
				} else if !reflect.DeepEqual(td.plaintext, res) {
					t.Fatalf("\nGot:  %s, \nWant: %s", bitstrings.RepresentPointerArray(res), bitstrings.RepresentPointerArray(td.plaintext))
				}
			},
		)
	}
}
