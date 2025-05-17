package encryption

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestMAC(t *testing.T) {
	tt := []struct {
		key                           string
		nonce                         *bitstrings.BitString128
		authenticatedData, ciphertext []*bitstrings.BitString128
		lengthAuth, lengthPlain       uint64
		mac                           *bitstrings.BitString128
	}{
		{
			key:   "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			nonce: bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
			authenticatedData: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("02 02 02 02 02 02 02 02 01 01 01 01 01 01 01 01"),
				bitstrings.FromGOSTString("04 04 04 04 04 04 04 04 03 03 03 03 03 03 03 03"),
				bitstrings.FromGOSTString("EA 05 05 05 05 05 05 05 05 00 00 00 00 00 00 00"),
			},
			ciphertext: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("A9 75 7B 81 47 95 6E 90 55 B8 A3 3D E8 9F 42 FC"),
				bitstrings.FromGOSTString("80 75 D2 21 2B F9 FD 5B D3 F7 06 9A AD C1 6b 39"),
				bitstrings.FromGOSTString("49 7A B1 59 15 A6 BA 85 93 6B 5D 0E A9 F6 85 1C"),
				bitstrings.FromGOSTString("C6 0C 14 D4 D3 F8 83 D0 AB 94 42 06 95 C7 6D EB"),
				bitstrings.FromGOSTString("2C 75 52 00 00 00 00 00 00 00 00 00 00 00 00 00"),
			},
			lengthAuth:  0x0000000000000148,
			lengthPlain: 0x0000000000000218,
			mac:         bitstrings.FromGOSTString("CF 5D 65 6F 40 C3 4F 5C 46 E8 BB 0E 29 FC DB 4C"),
		},
	}
	for i, td := range tt {
		t.Run(
			fmt.Sprintf("%d", i+1),
			func(t *testing.T) {
				keys, err := keyschedule.Schedule(td.key)
				if err != nil {
					t.Fatalf("error during keyschedule: %s", err)
				}
				res, err := computeMAC(
					keys,
					td.nonce,
					td.authenticatedData,
					td.ciphertext,
					td.lengthAuth,
					td.lengthPlain,
					context.Background(),
				)
				if err != nil {
					t.Fatalf("error during mac computation: %s", err)
				}
				if !reflect.DeepEqual(td.mac, res) {
					t.Fatalf("\nGot:  %0.16x, %0.16x,\nWant: %0.16x, %0.16x.", res.Upper(), res.Lower(), td.mac.Upper(), td.mac.Lower())
				}
			},
		)
	}
}
