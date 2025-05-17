package bitstrings_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestMSB(t *testing.T) {
	tt := []struct {
		in     *bitstrings.BitString128
		length int
		out    *bitstrings.BitString128
	}{
		{
			in:     bitstrings.FromGOSTString("2C 75 52 11 22 33 44 55 66 77 88 99 00 11 11 11"),
			length: 24,
			out:    bitstrings.FromGOSTString("2C 75 52 00 00 00 00 00 00 00 00 00 00 00 00 00"),
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%0.16x, %0.16x | %d -> %0.16x, %0.16x", td.in.Upper(), td.in.Lower(), td.length, td.out.Upper(), td.out.Lower()),
			func(t *testing.T) {
				res, err := bitstrings.MSB(td.in, td.length)
				if err != nil {
					t.Fatalf("Error: %s", err)
				}
				if res.Upper() != td.out.Upper() || res.Lower() != td.out.Lower() {
					t.Fatalf("\nGot:  %0.16x, %0.16x, \nWant: %0.16x, %0.16x.", res.Upper(), res.Lower(), td.out.Upper(), td.out.Lower())
				}
			},
		)
	}
}
