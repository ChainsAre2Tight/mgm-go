package gcm

import (
	"encoding/hex"
	"fmt"
	"reflect"
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
				keys, err := kuznechikgo.Schedule(k)
				if err != nil {
					t.Fatalf("Error during keyscheduling: %s", err)
				}
				bs, err := EncryptBitString(rawbs, keys)
				if err != nil {
					t.Fatalf("Error during encryption: %s", err)
				}
				if td.resU != bs.Upper() || td.resL != bs.Lower() {
					t.Fatalf("\nGot:  %x | %x, \nWant: %x | %x", bs.Upper(), bs.Lower(), td.resU, td.resL)
				}
			},
		)
	}
}

func TestEncryptBytes(t *testing.T) {
	tt := []struct {
		key     string
		in, out []byte
	}{
		{
			key: "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			in:  []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00, 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88},
			out: []byte{0x7F, 0x67, 0x9D, 0x90, 0xBE, 0xBC, 0x24, 0x30, 0x5A, 0x46, 0x8D, 0x42, 0xB9, 0xD4, 0xED, 0xCD},
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s | %v -> %v", td.key, td.in, td.out),
			func(t *testing.T) {
				k, err := hex.DecodeString(td.key)
				if err != nil {
					t.Fatalf("Error during key decoding: %s", err)
				}
				keys, err := kuznechikgo.Schedule(k)
				if err != nil {
					t.Fatalf("Error during keyschedule: %s", err)
				}
				if res, err := kuznechikgo.Encrypt(td.in, keys); err != nil {
					t.Fatalf("error during encryption: %s", err)
				} else if !reflect.DeepEqual(res, td.out) {
					t.Fatalf("\nGot:  %v, \nWant: %v.", res, td.out)
				}
			},
		)
	}
}
