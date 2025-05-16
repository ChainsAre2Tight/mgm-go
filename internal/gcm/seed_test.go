package gcm_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
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
			increment: bitstrings.IncremtntR,
			result: []*bitstrings.BitString128{
				bitstrings.FromGOSTString("B8 57 48 C5 12 F3 19 90 AA 56 7E F1 53 35 DB 74"),
				bitstrings.FromGOSTString("80 64 F0 12 6F AC 9B 2C 5B 6E AC 21 61 2F 94 33"),
				bitstrings.FromGOSTString("58 58 82 1D 40 C0 CD 0D 0A C1 E6 C2 47 09 8F 1C"),
				bitstrings.FromGOSTString("E4 3F 50 81 B5 8F 0B 49 01 2F 8E E8 6A CD 6D FA"),
				bitstrings.FromGOSTString("86 CE 9E 2A 0A 12 25 E3 33 56 91 B2 0D 5A 33 48"),
			},
		},
	}

	// TODO: fix prints
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s * %d + %s (%s) -> %v", td.iv.String(), td.depth, td.key, reflect.TypeOf(td.increment).Name(), td.result),
			func(t *testing.T) {
				keys, err := keyschedule.Schedule(td.key)
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

					t.Fatalf("\nGot:  %s, \nWant: %s.", representPointerArray(res), representPointerArray(td.result))
				}
			},
		)
	}
}

func representPointerArray(in []*bitstrings.BitString128) string {
	res := ""
	for _, pointer := range in {
		res = fmt.Sprintf("%s | %d, %d", res, pointer.Upper(), pointer.Lower())
	}
	return fmt.Sprintf("[%s]", res)
}
