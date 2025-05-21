package gcm_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/gcm"
)

func TestSeed(t *testing.T) {
	tt := []struct {
		iv        bitstrings.BitString128
		depth     int
		key       string
		increment func(*bitstrings.BitString128)
		result    []*bitstrings.BitString128
	}{
		{
			iv:        *bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
			depth:     5,
			key:       "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			increment: bitstrings.IncrementR,
			result: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("B8 57 48 C5 12 F3 19 90 AA 56 7E F1 53 35 DB 74"),
				bitstrings.FromGOSTString("80 64 F0 12 6F AC 9B 2C 5B 6E AC 21 61 2F 94 33"),
				bitstrings.FromGOSTString("58 58 82 1D 40 C0 CD 0D 0A C1 E6 C2 47 09 8F 1C"),
				bitstrings.FromGOSTString("E4 3F 50 81 B5 8F 0B 49 01 2F 8E E8 6A CD 6D FA"),
				bitstrings.FromGOSTString("86 CE 9E 2A 0A 12 25 E3 33 56 91 B2 0D 5A 33 48"),
			},
		}, {
			iv:        *bitstrings.FromGOSTString("91 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
			depth:     9,
			key:       "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF",
			increment: bitstrings.IncrementL,
			result: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("8D B1 87 D6 53 83 0E A4 BC 44 64 76 95 2C 30 0B"),
				bitstrings.FromGOSTString("7A 24 F7 26 30 E3 76 37 21 C8 F3 CD B1 DA 0E 31"),
				bitstrings.FromGOSTString("44 11 96 21 17 D2 06 35 C5 25 E0 A2 4D B4 B9 0A"),
				bitstrings.FromGOSTString("D8 C9 62 3C 4D BF E8 14 CE 7C 1C 0C EA A9 59 DB"),
				bitstrings.FromGOSTString("A5 E1 F1 95 33 3E 14 82 96 99 31 BF BE 6D FD 43"),
				bitstrings.FromGOSTString("B4 CA 80 8C AC CF B3 F9 17 24 E4 8A 2C 7E E9 D2"),
				bitstrings.FromGOSTString("72 90 8F C0 74 E4 69 E8 90 1B D1 88 EA 91 C3 31"),
				bitstrings.FromGOSTString("23 CA 27 15 b0 2C 68 31 3b FD ac b3 9E 4D 0F B8"),
				bitstrings.FromGOSTString("BC BC E6 C4 1A A3 55 A4 14 88 62 BF 64 BD 83 0D"),
			},
		},
	}

	// TODO: fix prints
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s * %d + %s (%s) -> %v", td.iv.String(), td.depth, td.key, reflect.TypeOf(td.increment).Name(), td.result),
			func(t *testing.T) {
				k, err := hex.DecodeString(td.key)
				if err != nil {
					t.Fatalf("Error during key decoding: %s", err)
				}
				keys, err := kuznechikgo.Schedule(k)
				if err != nil {
					t.Fatalf("Error during keyschedule: %s", err)
				}
				res, err := gcm.Seed(
					td.iv,
					td.depth,
					td.increment,
					keys,
					context.Background(),
				)
				if err != nil {
					t.Fatalf("Error during seeding: %s", err)
				}
				if !reflect.DeepEqual(res, td.result) {

					t.Fatalf("\nGot:  %s, \nWant: %s.", bitstrings.RepresentPointerArray(res), bitstrings.RepresentPointerArray(td.result))
				}
			},
		)
	}
}

func BenchmarkSeed(b *testing.B) {
	iv := *bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88")
	key, err := hex.DecodeString("8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF")
	if err != nil {
		b.Fatalf("Error during key decoding: %s", err)
	}
	keys, err := kuznechikgo.Schedule(key)
	if err != nil {
		b.Fatalf("Error during keyschedule: %s", err)
	}
	ctx := context.Background()
	for b.Loop() {
		_, err := gcm.Seed(
			iv,
			8,
			bitstrings.IncrementR,
			keys,
			ctx,
		)
		if err != nil {
			b.Fatalf("Error during seeding: %s", err)
		}
	}
}
