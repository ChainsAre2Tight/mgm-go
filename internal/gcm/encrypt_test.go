package gcm

import (
	"encoding/hex"
	"fmt"
	"testing"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestEncryptBitString(t *testing.T) {
	tt := []struct {
		key          string
		upper, lower uint64
		resU, resL   uint64
	}{
		{
			"8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			0x1122334455667700, 0xFFEEDDCCBBAA9988,
			0x7F679D90BEBC2430, 0x5A468D42B9D4EDCD,
		}, {
			"8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			0x9122334455667700, 0xFFEEDDCCBBAA9988,
			0x7FC245A8586E6602, 0xA7BBDB2786BDC66F,
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s + %d | %d ->  %d | %d", td.key, td.upper, td.lower, td.resU, td.resL),
			func(t *testing.T) {
				rawbs := bitstrings.NewBitString(
					td.upper,
					td.lower,
				)
				k, err := hex.DecodeString(td.key)
				if err != nil {
					t.Fatalf("Error during key decoding: %s", err)
				}
				fmt.Printf("%x\n", k)
				// fmt.Printf("%x, %x", rawbs.Upper(), rawbs.Lower())
				ks, err := kuznechikgo.Schedule(k)
				if err != nil {
					t.Fatalf("Error during keyscheduling: %s", err)
				}
				keys := kuznechikgo.KeysToUints(ks)
				bs := EncryptBitString(rawbs, keys)
				if td.resU != bs.Upper() || td.resL != bs.Lower() {
					t.Fatalf("\nGot:  %x | %x, \nWant: %x | %x", bs.Upper(), bs.Lower(), td.resU, td.resL)
				}
			},
		)
	}
}

func BenchmarkEncryptBitstring(b *testing.B) {
	key, err := hex.DecodeString("8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF")
	if err != nil {
		b.Fatalf("Error during key decoding: %s", err)
	}
	k, err := kuznechikgo.Schedule(key)
	if err != nil {
		b.Fatalf("Error during keyschedule: %s", err)
	}
	keys := kuznechikgo.KeysToUints(k)
	rawbs := bitstrings.NewBitString(0x1234567890123456, 0x1234567890123456)
	for b.Loop() {
		EncryptBitString(rawbs, keys)
	}
}
