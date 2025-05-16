package bitstrings_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestFromGOSTString(t *testing.T) {
	tt := []struct {
		in           string
		upper, lower uint64
	}{
		{"11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88", 1234605616436508416, 18441921395520346504},
		{"23 CA 27 15 B0 2C 68 31 3B FD AC B3 9E 4D 0F B8", 0x23CA2715B02C6831, 0x3BFDACB39E4D0FB8},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s -> %d | %d", td.in, td.lower, td.upper),
			func(t *testing.T) {
				res := bitstrings.FromGOSTString(td.in)
				if res == nil {
					t.Fatalf("nil res")
				}
				if res.Upper() != td.upper || res.Lower() != td.lower {
					t.Fatalf("\nGot:  %x, %x, \nWant: %x, %x.", res.Upper(), res.Lower(), td.upper, td.lower)
				}
			},
		)
	}
}
